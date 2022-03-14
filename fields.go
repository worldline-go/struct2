package struct2

import (
	"reflect"
)

func (d *Decoder) GetFields(s interface{}) []string {
	v := interface2StructValue(s)

	exportedFieldNames := []string{}

	getFields(v, d.tagName(), func(sf reflect.StructField) {
		if isFieldExported(sf) {
			tagName, _ := parseTag(sf.Tag.Get(d.tagName()))
			if tagName == "" {
				tagName = sf.Name
			}

			exportedFieldNames = append(exportedFieldNames, tagName)
		}
	})

	return exportedFieldNames
}

func isFieldExported(f reflect.StructField) bool {
	return f.PkgPath == ""
}

func getFields(v reflect.Value, tagName string, fn func(reflect.StructField)) {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if tag := field.Tag.Get(tagName); tag == "-" {
			continue
		}

		fn(field)
	}
}
