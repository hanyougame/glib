package errors

import (
	"errors"

	"go.temporal.io/sdk/temporal"
)

// Error 类型
const (
	TypeValidation    = "ValidationError"
	TypeNotFound      = "NotFoundError"
	TypeTimeout       = "TimeoutError"
	TypeDuplicate     = "DuplicateError"
	TypePermission    = "PermissionError"
	TypeConfiguration = "ConfigurationError"
	TypeExternal      = "ExternalServiceError"
	TypeGeneral       = "GeneralError"
)

// Error type 检查
func IsApplicationError(err error) bool {
	var appErr *temporal.ApplicationError
	return errors.As(err, &appErr)
}

func IsTimeoutError(err error) bool {
	var timeoutErr *temporal.TimeoutError
	return errors.As(err, &timeoutErr)
}

func IsCanceledError(err error) bool {
	var canceledErr *temporal.CanceledError
	return errors.As(err, &canceledErr)
}

func IsTerminatedError(err error) bool {
	var terminatedErr *temporal.TerminatedError
	return errors.As(err, &terminatedErr)
}

func IsActivityError(err error) bool {
	var activityErr *temporal.ActivityError
	return errors.As(err, &activityErr)
}

// Error creators
func NewApplicationError(message, errorType string, details ...interface{}) error {
	return temporal.NewApplicationError(message, errorType, details...)
}

func NewValidationError(message string, details ...interface{}) error {
	return temporal.NewApplicationError(message, TypeValidation, details...)
}

func NewNotFoundError(message string, details ...interface{}) error {
	return temporal.NewApplicationError(message, TypeNotFound, details...)
}

func NewTimeoutError(message string, details ...interface{}) error {
	return temporal.NewApplicationError(message, TypeTimeout, details...)
}

func NewDuplicateError(message string, details ...interface{}) error {
	return temporal.NewApplicationError(message, TypeDuplicate, details...)
}

func NewPermissionError(message string, details ...interface{}) error {
	return temporal.NewApplicationError(message, TypePermission, details...)
}

func NewConfigurationError(message string, details ...interface{}) error {
	return temporal.NewApplicationError(message, TypeConfiguration, details...)
}

func NewExternalServiceError(message string, details ...interface{}) error {
	return temporal.NewApplicationError(message, TypeExternal, details...)
}
