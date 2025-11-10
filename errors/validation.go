package errors

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// ValidateNotNil checks if parameters are nil and returns an error if any is nil
func ValidateNotNil(parameters map[string]interface{}) error {
	var nilParams []string

	for name, parameter := range parameters {
		if parameter == nil {
			nilParams = append(nilParams, name)
			continue
		}

		val := reflect.ValueOf(parameter)
		switch val.Kind() {
		case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
			if val.IsNil() {
				nilParams = append(nilParams, name)
			}
		}
	}

	if len(nilParams) > 0 {
		pc := make([]uintptr, 1)
		runtime.Callers(2, pc)
		funcName := runtime.FuncForPC(pc[0]).Name()
		shortName := funcName[strings.LastIndexByte(funcName, '/')+1:]

		return fmt.Errorf("func `%s`: parameters cannot be nil: %s",
			shortName,
			strings.Join(nilParams, ", "))
	}

	return nil
}
