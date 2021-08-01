package errors

import "reflect"

// CustomError is an customized error
type CustomError struct {
	Message       string
	InternalError error
}

// Error shows error message
func (c *CustomError) Error() string {
	return c.Message
}

// InternalErrorTypeString shows error message
func (c *CustomError) InternalErrorTypeString() string {
	return reflect.Indirect(reflect.ValueOf(c.InternalError)).Type().String()
}
