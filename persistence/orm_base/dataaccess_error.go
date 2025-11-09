package orm_base
import (
	"github.com/janmbaco/go-infrastructure/errors"
)

type DataBaseError interface {
	errors.CustomError
	GetErrorType() DataBaseErrorType
}

type databaseError struct {
	errors.CustomizableError
	ErrorType DataBaseErrorType
}

func newDataBaseError(errorType DataBaseErrorType, message string, internalError error) DataBaseError {
	return &databaseError{
		CustomizableError: errors.CustomizableError{
			Message:       message,
			InternalError: internalError,
		},
		ErrorType: errorType,
	}
}

func (e *databaseError) GetErrorType() DataBaseErrorType {
	return e.ErrorType
}

// DataBaseErrorType is the type of the errors from database
type DataBaseErrorType uint8

const (
	UnexpectedError DataBaseErrorType = iota
	DataRowUnexpected
	DataFilterUnexpected
)


