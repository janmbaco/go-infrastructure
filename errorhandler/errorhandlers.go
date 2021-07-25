package errorhandler

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func OnErrorContinue(callBack func()) {
	defer deferTryError(true, nil, nil)
	callBack()
}

func TryCatchError(callBack func(), errorFunc func(error)) {
	defer deferTryError(false, errorFunc, nil)
	callBack()
}

func TryFinally(callBack func(), finallyFunc func()) {
	defer deferTryError(false, nil, finallyFunc)
	callBack()
}

func TryCatchErrorAndFinally(callBack func(), errorFunc func(error), finallyFunc func()) {
	defer deferTryError(false, errorFunc, finallyFunc)
	callBack()
}

func deferTryError(shouldContinue bool, errorFunc func(error), finallyFunc func()) {
	if finallyFunc != nil {
		finallyFunc()
	}
	if re := recover(); re != nil {
		text := "unexpected error"
		switch re.(type) {
		case string:
			text = re.(string)
		case error:
			text = re.(error).Error()
		}
		if errorFunc != nil {
			errorFunc(errors.New(text))
		} else if !shouldContinue {
			panic(errors.New(text))
		}

	}
}

func TryPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func TryPanicError(err error, errorFunc func(error)) {
	defer deferTryError(false, errorFunc, nil)
	TryPanic(err)
}

func TryPanicFinally(err error, finallyFunc func()) {
	defer deferTryError(false, nil, finallyFunc)
	TryPanic(err)
}

func TryPanicErrorAndFinally(err error, errorFunc func(error), finallyFunc func()) {
	defer deferTryError(false, errorFunc, finallyFunc)
	TryPanic(err)
}

func CheckNilParameter(parameters map[string]interface{}) {
	panicMessage := make([]string, 1)
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	funcName := runtime.FuncForPC(pc[0]).Name()
	panicMessage[0] = fmt.Sprintf("Func `%v`", funcName[strings.LastIndexByte(funcName, '/')+1:])
	for name, parameter := range parameters {
		if typ := reflect.TypeOf(parameter); typ.Kind() == reflect.Ptr && reflect.ValueOf(parameter).IsNil() {
			panicMessage = append(panicMessage, fmt.Sprintf("The parameter '%v' can't be nil!", name))
		}
	}
	if len(panicMessage) > 1 {
		panic(strings.Join(panicMessage, "\n"))
	}
}
