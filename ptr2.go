package struct2

import "reflect"

func Ptr2Concrete(val any) any {
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Pointer {
		if !v.IsNil() {
			return v.Elem().Interface()
		}

		// create new value from that type
		return reflect.Zero(v.Type().Elem()).Interface()
	}

	return val
}
