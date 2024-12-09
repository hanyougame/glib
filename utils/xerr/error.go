package xerr

import "github.com/zeromicro/x/errors"

// New 自定义错误
func New(code ErrCode, msg string) error {
	return errors.New(int(code), msg)
}

// NewParamErr 参数错误
func NewParamErr(msg string) error {
	return New(ParamError, msg)
}

// NewServerErr server error
func NewServerErr(msg string) error {
	return New(ServerError, msg)
}

// NewDbErr db error
func NewDbErr(msg string) error {
	return New(DbError, msg)
}
