package test

//go:generate go-easy generate enum --type EnumTest --type EnumTestB
type EnumTest uint64
type EnumTestB uint64

const (
	// 占位
	EnumTestNone EnumTest = iota
	// 测试一号
	EnumTest1
	// 测试二号
	EnumTest2
	// 占位
	EnumTestBNone EnumTestB = iota
	// 测试一号
	EnumTestB1
	// 测试二号
	EnumTestB2
)
