package ginx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

type Response struct {
	ErrorCode uint64      `json:"error_code"`
	Data      interface{} `json:"data"`
	Message   string      `json:"message"`
}

type ErrorCode interface {
	error
	fmt.Stringer
	ErrorCode() uint64
}

var (
	InternalError ErrorCode = nil
	initCheckOnce           = &sync.Once{}
)

func initCheck() {
	initCheckOnce.Do(func() {
		if InternalError == nil {
			panic("InternalError is nil")
		}
	})
}

func SetInternalError(value ErrorCode) {
	InternalError = value
}

func MakeDataResponse(c *gin.Context, data interface{}) {
	initCheck()
	c.JSON(http.StatusOK, Response{Data: data})
}

func MakeMessageResponse(c *gin.Context, message string) {
	initCheck()
	c.JSON(http.StatusOK, Response{Message: message})
}

func MakeErrorResponse(c *gin.Context, err error) {
	initCheck()
	resp := Response{
		Message: err.Error(),
	}
	if errorCode, ok := err.(ErrorCode); ok {
		resp.ErrorCode = errorCode.ErrorCode()
	} else {
		resp.ErrorCode = InternalError.ErrorCode()
	}
	c.JSON(http.StatusOK, resp)
}
