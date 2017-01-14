package utils

import (
	"fmt"
	"reflect"
)

func IsPrimitive(data interface{}) (ok bool) {
	if data == nil {
		return false
	}
	var kind reflect.Kind

	if kind, ok = data.(reflect.Kind); !ok {
		kind = reflect.TypeOf(data).Kind()
	}

	switch kind {
	case reflect.Bool:
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.String:
		return true
	}

	return false
}

func ToString(args ...interface{}) string {
	result := ""
	for _, arg := range args {
		if arg != nil {
			result += fmt.Sprintf("%v", arg)
		}
	}
	return result
}
