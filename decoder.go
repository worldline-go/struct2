package struct2

import (
	"reflect"
)

// Decoder is main struct of struct2, holds config and functions.
type Decoder struct {
	// Tagname to lookup struct's field tag.
	TagName string // default is 'struct'
	// Hooks function run before decode and enable to change of data.
	Hooks []HookFunc
}

func (d *Decoder) tagName() string {
	if d.TagName == "" {
		return "struct"
	}

	return d.TagName
}

func interface2StructValue(s interface{}) reflect.Value {
	v := reflect.ValueOf(s)

	// pointer to struct
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		panic("not struct value")
	}

	return v
}
