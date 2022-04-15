package errorstest

import (
	"testing"

	"github.com/janmbaco/go-infrastructure/dependencyinjection/static"
	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	errorsResolver "github.com/janmbaco/go-infrastructure/errors/ioc/resolver"
	"github.com/janmbaco/go-infrastructure/logs"
)

func TestErrorManager(t *testing.T) {

	static.Container.Register().AsType(new(sampleClass), newSampleCass, nil)
	log := static.Container.Resolver().Type(new(logs.Logger), nil).(logs.Logger)
	errorsResolver.GetErrorManager().On(
		new(sampleClassError), 
		func(err error) {
			switch err.(sampleClassError).GetErrorType() {
			case _UnexpectedError:
				log.Info("UnexpectedError")
			case _Error1:
				log.Info("Error1")
			case _Error2:
				log.Info("Error2")
				
			}
		})

		sampleClass := static.Container.Resolver().Type(new(sampleClass), nil).(sampleClass)

		sampleClass.PanicUnexpected()

		sampleClass.PanicError1()

		sampleClass.PanicError2()
}

type sampleClassError interface {
	errors.CustomError
	GetErrorType() sampleClassErrorType
}

type sampleClassErrorType uint8

const (
	_UnexpectedError sampleClassErrorType = iota
	_Error1
	_Error2
)

type sampleClassErrorImp struct{
	errors.CustomizableError
	ErrorType sampleClassErrorType
}

func(s *sampleClassErrorImp) GetErrorType() sampleClassErrorType{
	return s.ErrorType
}

type sampleClass interface{
	PanicUnexpected()
	PanicError1()
	PanicError2()
}

type sampleClassImp struct
{ 
	errorDefer errors.ErrorDefer
}

func newSampleCass(errorDefer errors.ErrorDefer) sampleClass{
	errorschecker.CheckNilParameter(map[string]interface{}{"errorDefer":errorDefer})
	return &sampleClassImp{errorDefer: errorDefer}
}

func (s *sampleClassImp) PanicUnexpected(){
	defer s.errorDefer.TryThrowError(nil) 
	panic(&sampleClassErrorImp{ ErrorType: _UnexpectedError})
 }

func (s *sampleClassImp) PanicError1(){
	defer s.errorDefer.TryThrowError(nil)
	panic(&sampleClassErrorImp{ ErrorType: _Error1})
}

func (s *sampleClassImp) PanicError2(){
	defer s.errorDefer.TryThrowError(nil)
	panic(&sampleClassErrorImp{ ErrorType: _Error2})
}