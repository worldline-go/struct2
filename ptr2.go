package struct2

import "reflect"

func Ptr2Concrete(val interface{}) interface{} {
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr {
		if !v.IsNil() {
			return v.Elem().Interface()
		}

		// create new value from that type
		return reflect.Zero(v.Type().Elem()).Interface()
	}

	return val
}
