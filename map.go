package struct2

import (
	"reflect"
	"strings"
)

type configMap struct {
	omitNested bool
}

type optionMap func(*configMap)

func WithOmitNested() optionMap {
	return func(c *configMap) {
		c.omitNested = true
	}
}

// MapOmitNested converts given struct to the map[string]any omitting nested structs and maps.
//
// Deprecated: use Map with WithOmitNested option instead.
func (d *Decoder) MapOmitNested(input any) map[string]any {
	return d.Map(input, WithOmitNested())
}

// MapSlice converts given slice of structs to []map[string]any.
// Panic if input not a slice or array type.
func (d *Decoder) MapSlice(input any, opts ...optionMap) []map[string]any {
	config := configMap{}
	for _, opt := range opts {
		opt(&config)
	}

	inputV := reflect.ValueOf(input)
	if isNil(inputV) {
		return nil
	}

	if inputV.Kind() != reflect.Slice && inputV.Kind() != reflect.Array {
		panic("input not a slice or array type")
	}

	out := make([]map[string]any, inputV.Len())
	for i := 0; i < inputV.Len(); i++ {
		out[i] = d.convertMap(inputV.Index(i).Interface(), config)
	}

	return out
}

// Map converts given struct to the map[string]any.
// Panic if input not a struct type.
func (d *Decoder) Map(input any, opts ...optionMap) map[string]any {
	config := configMap{}
	for _, opt := range opts {
		opt(&config)
	}

	return d.convertMap(input, config)
}

func (d *Decoder) convertMap(input any, config configMap) map[string]any {
	inputV := reflect.ValueOf(input)
	if isNil(inputV) {
		return nil
	}

	v := value2StructValue(inputV)

	out := make(map[string]any)

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
		var finalVal any

		tagName, tagOpts := d.parseTag(field)
		if tagName != "" {
			name = tagName
		}

		// if the value is a zero value and the field is marked as omitempty do
		// not include
		if tagOpts.Has("omitempty") && val.IsZero() {
			continue
		}

		if d.OmitNilPtr && val.Kind() == reflect.Pointer && val.IsNil() {
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
			if val.Type().Kind() == reflect.Pointer && val.IsNil() {
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
			if v.Kind() == reflect.Pointer {
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
			for k := range finalVal.(map[string]any) {
				out[k] = finalVal.(map[string]any)[k]
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
func (d *Decoder) nested(val reflect.Value) any {
	var finalVal any

	v := reflect.ValueOf(val.Interface())
	if v.Kind() == reflect.Pointer {
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
		case reflect.Pointer, reflect.Array, reflect.Map,
			reflect.Slice, reflect.Chan:
			mapElem = val.Type().Elem()
			if mapElem.Kind() == reflect.Pointer {
				mapElem = mapElem.Elem()
			}
		}

		// only iterate over struct types, ie: map[string]StructType,
		// map[string][]StructType,
		if mapElem.Kind() == reflect.Struct ||
			(mapElem.Kind() == reflect.Slice && mapElem.Elem().Kind() == reflect.Struct) {
			m := make(map[string]any, val.Len())
			for _, k := range val.MapKeys() {
				m[k.String()] = d.nested(val.MapIndex(k))
			}

			finalVal = m

			break
		}

		// TODO(arslan): should this be optional?
		finalVal = val.Interface()
	case reflect.Slice, reflect.Array:
		if val.Type().Kind() == reflect.Pointer {
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
			!(val.Type().Elem().Kind() == reflect.Pointer && val.Type().Elem().Elem().Kind() == reflect.Struct) {
			finalVal = val.Interface()

			break
		}

		slices := make([]any, val.Len())
		for x := 0; x < val.Len(); x++ {
			slices[x] = d.nested(val.Index(x))
		}

		finalVal = slices
	default:
		finalVal = val.Interface()
	}

	return finalVal
}
