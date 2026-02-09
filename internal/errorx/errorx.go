package errorx

const DefaultCode = 1001

type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type CodeErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewCodeError(code int, msg string) error {
	return &CodeError{Code: code, Msg: msg}
}

func NewDefaultError(msg string) error {
	return NewCodeError(DefaultCode, msg)
}

func (e *CodeError) Error() string {
	return e.Msg
}

func (e *CodeError) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		Code: e.Code,
		Msg:  e.Msg,
	}
}

// 预定义业务错误码
const (
	ErrCodeOK             = 200
	ErrCodeParamInvalid   = 400
	ErrCodeUnauthorized   = 401
	ErrCodeForbidden      = 403
	ErrCodeNotFound       = 404
	ErrCodeServerInternal = 500

	// 业务相关
	ErrCodeUserNotFound     = 10001
	ErrCodeUserAlreadyExist = 10002
	ErrCodeCacheError       = 10003
)
