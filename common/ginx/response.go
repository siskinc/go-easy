package ginx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Result interface{} `json:"result"`
}

type ErrorResponse struct {
	ErrorCode uint64 `json:"error_code"`
	Message   string `json:"message"`
}

type ErrorCode interface {
	error
	fmt.Stringer
	ErrorCode() uint64
}

func MakeResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Result: data})
}

func MakeErrorResponse(c *gin.Context, err error) {
	resp := ErrorResponse{
		Message: err.Error(),
	}
	if errorCode, ok := err.(ErrorCode); ok {
		resp.ErrorCode = errorCode.ErrorCode()
		c.JSON(HttpStatusCustomInternalServerError, resp)
	} else {
		c.JSON(http.StatusBadRequest, resp)
	}
}
