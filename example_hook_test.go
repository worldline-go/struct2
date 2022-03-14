package struct2_test

import (
	"fmt"
	"reflect"

	"github.com/worldline-go/struct2"
)

func Example_customHook() {
	type ColorGroup struct {
		Name  string `db:"name"`
		Count int    `db:"count"`
	}

	group := ColorGroup{
		Name: "DeepCore",
	}

	decoder := struct2.Decoder{
		TagName: "db",
		Hooks: []struct2.HookFunc{func(v reflect.Value) (interface{}, error) {
			if v.Kind() == reflect.String {
				return "str_" + v.Interface().(string), nil
			}

			return nil, struct2.ErrContinueHook
		}},
	}

	result := decoder.Map(group)

	fmt.Printf("%v", result["name"])
	// Output:
	// str_DeepCore
}
