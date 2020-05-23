package errorhandler

import (
	"errors"
	"github.com/janmbaco/go-infrastructure/logs"
)

func TryError(callBack func(), errorFunc func(error)) {
	func() {
		defer func() {
			if re := recover(); re != nil {
				text := "unexpected error"
				switch re.(type) {
				case string:
					text = re.(string)
				case error:
					text = re.(error).Error()
				}
				errorFunc(errors.New(text))
			}
		}()
		callBack()
	}()
}

func TryFinally(callBack func(), finallyFunc func()) {
	func() {
		defer func() {
			finallyFunc()
			if re := recover(); re != nil {
				text := "unexpected error"
				switch re.(type) {
				case string:
					text = re.(string)
				case error:
					text = re.(error).Error()
				}
				panic(errors.New(text))
			}
		}()
		callBack()
	}()
}

func TryPanic(err error) {
	if err != nil {
		logs.Log.Error(err.Error())
		panic(err)
	}
}
