package ginx

import "github.com/gin-gonic/gin"

func HttpRegister(router *gin.RouterGroup, relativePath string, restAPI interface{}, handlers ...gin.HandlerFunc) {
	if get, ok := restAPI.(Get); ok {
		httpGetPath := ""
		if getPath, ok := restAPI.(HttpGetPathDescriptor); ok {
			httpGetPath = getPath.HttpGetPath()
		}
		router.GET(MergeURL(relativePath, httpGetPath), get.GET)
	}

	if post, ok := restAPI.(Post); ok {
		httpPostPath := ""
		if postPath, ok := restAPI.(HttpPostPathDescriptor); ok {
			httpPostPath = postPath.HttpPostPath()
		}
		router.POST(MergeURL(relativePath, httpPostPath), post.POST)
	}

	if put, ok := restAPI.(Put); ok {
		httpPutPath := ""
		if putPath, ok := restAPI.(HttpPutPathDescriptor); ok {
			httpPutPath = putPath.HttpPutPath()
		}
		router.PUT(MergeURL(relativePath, httpPutPath), put.PUT)
	}

	if patch, ok := restAPI.(Patch); ok {
		httpPatchPath := ""
		if patchPath, ok := restAPI.(HttpPatchPathDescriptor); ok {
			httpPatchPath = patchPath.HttpPatchPath()
		}
		router.PATCH(MergeURL(relativePath, httpPatchPath), patch.PATCH)
	}

	if delete_, ok := restAPI.(Delete); ok {
		httpDeletePath := ""
		if deletePath, ok := restAPI.(HttpDeletePathDescriptor); ok {
			httpDeletePath = deletePath.HttpDeletePath()
		}
		router.DELETE(MergeURL(relativePath, httpDeletePath), delete_.DELETE)
	}
}
