package errors

import (
	"fmt"
)

type CodeError struct {
	errCode int
	errMsg  string
}

func New(errCode int, errMsg string) *CodeError {
	return &CodeError{errCode: errCode, errMsg: errMsg}
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode:%d，ErrMsg:%s", e.errCode, e.errMsg)
}

func (e *CodeError) ErrCode() int {
	return e.errCode
}

func (e *CodeError) ErrMsg() string {
	return e.errMsg
}
