package errorhandler

import (
	"errors"

	"github.com/janmbaco/go-infrastructure/logs"
)

func onErrorContinue(callBack func()) {
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
		logs.Log.Error(text)
		if errorFunc != nil {
			errorFunc(errors.New(text))
		} else if !shouldContinue {
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
