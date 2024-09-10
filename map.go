package struct2

import (
	"reflect"
	"strings"
)

type configMap struct {
	omitNested bool
}

// MapOmitNested converts given struct to the map[string]interface{} without looking nested object.
// Panic if input not a struct type.
func (d *Decoder) MapOmitNested(input interface{}) map[string]interface{} {
	return d.convertMap(input, configMap{omitNested: true})
}

// Map converts given struct to the map[string]interface{}.
// Panic if input not a struct type.
func (d *Decoder) Map(input interface{}) map[string]interface{} {
	return d.convertMap(input, configMap{})
}

func (d *Decoder) convertMap(input interface{}, config configMap) map[string]interface{} {
	inputV := reflect.ValueOf(input)
	if isNil(inputV) {
		return nil
	}

	v := value2StructValue(inputV)

	out := make(map[string]interface{})

	var fields []reflect.StructField

	d.getFields(v, func(sf reflect.StructField) {
		fields = append(fields, sf)
	})

FIELDS:
	for _, field := range fields {
		name := field.Name
		val := v.FieldByName(name)
		if d.OuputCamelCase {
			name = strings.ToLower(name[0:1]) + name[1:]
		}

		isSubStruct := false
		var finalVal interface{}

		tagName, tagOpts := d.parseTag(field)
		if tagName != "" {
			name = tagName
		}

		// if the value is a zero value and the field is marked as omitempty do
		// not include
		if tagOpts.Has("omitempty") && val.IsZero() {
			continue
		}

		if d.OmitNilPtr && val.Kind() == reflect.Ptr && val.IsNil() {
			continue
		}

		if tagOpts.Has("string") {
			s, err := toStringE(val.Interface())
			if err != nil {
				continue
			}
			out[name] = s
			continue
		}

		ptr2 := d.ForcePtr2 || tagOpts.Has("ptr2")

		// custom hooks
		for _, hook := range d.Hooks {
			if hookResult, err := hook(val); err == nil {
				if ptr2 {
					out[name] = Ptr2Concrete(hookResult)
				} else {
					out[name] = hookResult
				}

				continue FIELDS
			}
		}

		// type hook
		var hook Hooker
		if hookSelect, ok := val.Interface().(Hooker); ok {
			hook = hookSelect
		} else {
			addrVal := reflect.New(val.Type())
			reflect.Indirect(addrVal).Set(val)

			if hookSelect, ok := addrVal.Interface().(Hooker); ok {
				hook = hookSelect
			}
		}

		if hook != nil {
			if val.Type().Kind() == reflect.Ptr && val.IsNil() {
				// nil pointer call to value receiver
				out[name] = nil

				continue
			}

			if ptr2 {
				out[name] = Ptr2Concrete(hook.Struct2Hook())
			} else {
				out[name] = hook.Struct2Hook()
			}

			continue
		}

		// nested parts

		if !config.omitNested && !tagOpts.Has("omitnested") {
			finalVal = d.nested(val)

			v := reflect.ValueOf(val.Interface())
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}

			switch v.Kind() {
			case reflect.Map, reflect.Struct:
				isSubStruct = true
			}
		} else {
			finalVal = val.Interface()
		}

		if isSubStruct && (tagOpts.Has("flatten")) {
			for k := range finalVal.(map[string]interface{}) {
				out[k] = finalVal.(map[string]interface{})[k]
			}
		} else {
			if ptr2 {
				out[name] = Ptr2Concrete(finalVal)
			} else {
				out[name] = finalVal
			}
		}
	}

	return out
}

// nested retrieves recursively all types for the given value and returns the nested value.
func (d *Decoder) nested(val reflect.Value) interface{} {
	var finalVal interface{}

	v := reflect.ValueOf(val.Interface())
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		exportedFieldCount := 0

		d.getFields(v, func(sf reflect.StructField) {
			if isFieldExported(sf) {
				exportedFieldCount++
			}
		})

		if exportedFieldCount > 0 {
			finalVal = d.Map(val.Interface())
		} else {
			finalVal = val.Interface()
		}
	case reflect.Map:
		// get the element type of the map
		mapElem := val.Type()

		switch val.Type().Kind() {
		case reflect.Ptr, reflect.Array, reflect.Map,
			reflect.Slice, reflect.Chan:
			mapElem = val.Type().Elem()
			if mapElem.Kind() == reflect.Ptr {
				mapElem = mapElem.Elem()
			}
		}

		// only iterate over struct types, ie: map[string]StructType,
		// map[string][]StructType,
		if mapElem.Kind() == reflect.Struct ||
			(mapElem.Kind() == reflect.Slice && mapElem.Elem().Kind() == reflect.Struct) {
			m := make(map[string]interface{}, val.Len())
			for _, k := range val.MapKeys() {
				m[k.String()] = d.nested(val.MapIndex(k))
			}

			finalVal = m

			break
		}

		// TODO(arslan): should this be optional?
		finalVal = val.Interface()
	case reflect.Slice, reflect.Array:
		if val.Type().Kind() == reflect.Ptr {
			val = val.Elem()
		}

		if val.Type().Kind() == reflect.Interface {
			finalVal = val.Interface()

			break
		}

		// TODO(arslan): should this be optional?
		// do not iterate of non struct types, just pass the value. Ie: []int,
		// []string, co... We only iterate further if it's a struct.
		// i.e []foo or []*foo
		if val.Type().Elem().Kind() != reflect.Struct &&
			!(val.Type().Elem().Kind() == reflect.Ptr && val.Type().Elem().Elem().Kind() == reflect.Struct) {
			finalVal = val.Interface()

			break
		}

		slices := make([]interface{}, val.Len())
		for x := 0; x < val.Len(); x++ {
			slices[x] = d.nested(val.Index(x))
		}

		finalVal = slices
	default:
		finalVal = val.Interface()
	}

	return finalVal
}
