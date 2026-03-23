package handler

// Response 标准返回结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Success 成功返回的快捷工具
func Success(data interface{}) Response {
	return Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	}
}

// Error 失败返回的快捷工具
func Error(code int, msg string) Response {
	return Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}
