package xerr

// ErrCode 业务错误码
type ErrCode int

// Int ErrCode转int
func (e ErrCode) Int() int {
	return int(e)
}

const (
	ParamError        ErrCode = 400 // 参数错误
	UnauthorizedError ErrCode = 401 // 无权限
	ServerError       ErrCode = 500 // 服务器内部错误
	DbError           ErrCode = 600 // 数据库错误
	CaptchaError      ErrCode = 700 // 验证码错误
)

var (
	ErrParam        = New(ParamError, "param error")
	ErrUnauthorized = New(UnauthorizedError, "unauthorized error")
	ErrServer       = New(ServerError, "server error")
	ErrDB           = New(DbError, "db error")
	ErrCaptcha      = New(CaptchaError, "captcha error")
)
