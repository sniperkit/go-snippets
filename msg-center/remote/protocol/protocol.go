package protocol

const (
	// GUnknownProtoCode 这个代表着发生了未知的协议
	GUnknownProtoCode = 0

	// TestProtoCode 是一个测试协议
	TestProtoCode = 1001
)

//  请求协议
type (
	// TestProtoRequest 响应请求协议
	TestProtoRequest struct {
		Name    string
		Age     string
		Address string
	}
)

// 响应协议
type ()
