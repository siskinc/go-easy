package ginx

import (
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

func MakeResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Result: data})
}

func MakeErrorResponse(c *gin.Context, errorCode uint64, message string) {
	c.JSON(HttpStatusCustomInternalServerError, ErrorResponse{
		ErrorCode: errorCode,
		Message:   message,
	})
}
