package struct2

import (
	"reflect"
)

func (d *Decoder) GetFields(s interface{}) []string {
	v := value2StructValue(reflect.ValueOf(s))

	exportedFieldNames := []string{}

	d.getFields(v, func(sf reflect.StructField) {
		if isFieldExported(sf) {

			tagName, _ := d.parseTag(sf)
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

func (d *Decoder) getFields(v reflect.Value, fn func(reflect.StructField)) {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		tagValue, _ := d.getTagValue(field)
		if tag := tagValue; tag == "-" {
			continue
		}

		fn(field)
	}
}
