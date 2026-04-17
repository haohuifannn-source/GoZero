package errorx

const (
	defaultErrCode = 1001
	RPCErrCode     = 1002
)

// CodeErr自定义错误
type CodeErr struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// CodeErrorResponse自定义返回错误
type CodeErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// NewCodeError自定义错误返回
func NewCodeError(code int, msg string) error {
	return &CodeErr{
		Code: code,
		Msg:  msg,
	}
}

// NewDefaultCodeError返回默认的错误
func NewDefaultCodeError(msg string) error {
	return &CodeErr{
		Code: defaultErrCode,
		Msg:  msg,
	}
}

// CodeErr实现了error接口-----注意这里是指针实现，需要在返回的时候用&
func (e *CodeErr) Error() string {
	return e.Msg
}

// Data返回自定义错误
func (e *CodeErr) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		Code: e.Code,
		Msg:  e.Msg,
	}
}
