package ginx

type HttpGetPathDescriptor interface {
	HttpGetPath() string
}

type HttpPostPathDescriptor interface {
	HttpPostPath() string
}

type HttpPutPathDescriptor interface {
	HttpPutPath() string
}

type HttpPatchPathDescriptor interface {
	HttpPatchPath() string
}

type HttpDeletePathDescriptor interface {
	HttpDeletePath() string
}
