package xhttp

// ErrorResult is a struct that contains a code and a message.
// It implements the error interface.
type ErrorResult struct {
	Code    int
	Message string
	Data    any
}

func (c *ErrorResult) Error() string {
	return c.Message
}

// New creates a new CodeMsg.
func New(code int, msg string, data any) error {
	return &ErrorResult{Code: code, Message: msg, Data: data}
}
