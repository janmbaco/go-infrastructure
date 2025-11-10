package errors

// CustomError is an customized error
type CustomError interface {
	error
	GetMessage() string
	GetInternalError() error
}

// CustomizableError is a error with a error with a customizable message
type CustomizableError struct {
	InternalError error
	Message       string
}

func (e *CustomizableError) Error() string {
	return e.Message
}

func (e *CustomizableError) GetMessage() string {
	return e.Message
}

func (e *CustomizableError) GetInternalError() error {
	return e.InternalError
}
