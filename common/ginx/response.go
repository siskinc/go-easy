package ginx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	ErrorCode uint64 `json:"error_code"`
}

type DataResponse struct {
	Response
	Data interface{} `json:"data"`
}

type MessageResponse struct {
	Response
	Message string `json:"message"`
}

type ErrorCode interface {
	error
	fmt.Stringer
	ErrorCode() uint64
}

func MakeDataResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, DataResponse{Data: data})
}

func MakeMessageResponse(c *gin.Context, message string) {
	c.JSON(http.StatusOK, MessageResponse{Message: message})
}

func MakeErrorResponse(c *gin.Context, err error) {
	resp := MessageResponse{
		Message: err.Error(),
	}
	if errorCode, ok := err.(ErrorCode); ok {
		resp.ErrorCode = errorCode.ErrorCode()
		c.JSON(http.StatusOK, resp)
	} else {
		resp.ErrorCode = 666666
	}
}
