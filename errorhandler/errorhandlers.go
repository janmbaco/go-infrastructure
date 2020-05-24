package errorhandler

import (
	"errors"

	"github.com/janmbaco/go-infrastructure/logs"
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
			logs.Log.Error(text)
			panic(errors.New(text))
		}

	}
}

func TryPanic(err error) {
	if err != nil {
		logs.Log.Error(err.Error())
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
