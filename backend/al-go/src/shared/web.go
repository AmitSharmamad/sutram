package shared

/// 4xx status codes

type BadRequestError struct {
	Message string
	Code    string
}

func NewBadRequest(msg string, code string) *BadRequestError {
	return &BadRequestError{
		Message: msg,
		Code:    code,
	}
}

func (e *BadRequestError) Error() string {
	return e.Message
}

/// 5xx status codes
