package ginx

import "github.com/gin-gonic/gin"

func BindAll(c *gin.Context, data interface{}) (err error) {
	if err = c.BindUri(data); nil != err {
		return
	}

	if err = c.BindQuery(data); nil != err {
		return
	}

	if err = c.BindYAML(data); nil != err {
		return
	}

	if err = c.BindJSON(data); nil != err {
		return
	}

	return
}
