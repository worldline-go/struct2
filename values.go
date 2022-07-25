package struct2

import (
	"reflect"
)

// isNil return true if the given value is nil.
//
//   isNil(reflect.ValueOf(x))
func isNil(v reflect.Value) bool {
	if v.Kind() == reflect.Ptr {
		return v.IsNil()
	}

	return false
}

func value2StructValue(v reflect.Value) reflect.Value {
	// pointer to struct
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.Kind() == reflect.Interface {
			v = reflect.Indirect(v)
		}

		if v.IsNil() {
			v = reflect.New(v.Type().Elem()).Elem()
			// v = reflect.Zero(v.Type().Elem())

			break
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		panic("not struct value")
	}

	return v
}
