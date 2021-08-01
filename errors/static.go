package errors

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// CheckNilParameter checks if the parameters are nil
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

// TryPanic throws an error to panic if error is different from nil
func TryPanic(err error) {
	if err != nil {
		panic(err)
	}
}
