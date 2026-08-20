package pkg

// Response 统一返回结果（兼容前端 Result<T> 格式）
type Response struct {
	Code int         `json:"code"` // 1-成功 0-失败
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// PageResult 分页结果
type PageResult struct {
	Total   int64       `json:"total"`
	Records interface{} `json:"records"`
}

func Success(data interface{}) Response {
	return Response{Code: 1, Data: data}
}

func SuccessMsg(msg string) Response {
	return Response{Code: 1, Msg: msg}
}

func Error(msg string) Response {
	return Response{Code: 0, Msg: msg}
}

func ErrorWithData(msg string, data interface{}) Response {
	return Response{Code: 0, Msg: msg, Data: data}
}
