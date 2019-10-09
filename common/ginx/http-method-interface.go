package ginx

import "github.com/gin-gonic/gin"

type Get interface {
	GET(c *gin.Context)
}

type Post interface {
	POST(c *gin.Context)
}

type Put interface {
	PUT(c *gin.Context)
}

type Patch interface {
	PATCH(c *gin.Context)
}

type Delete interface {
	DELETE(c *gin.Context)
}
