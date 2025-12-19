package struct2_test

import (
	"fmt"
	"reflect"
	"time"

	"github.com/worldline-go/struct2"
)

func Example_mapToStruct() {
	type SubCfg struct {
		Enabled bool `cfg:"enabled"`
	}

	type Config struct {
		Name  string        `cfg:"name"`
		Count int           `cfg:"count"`
		Sub   SubCfg        `cfg:"sub"`
		Time  time.Duration `cfg:"time"`
	}

	decoder := struct2.Decoder{
		TagName: "cfg",
		HooksDecode: []struct2.HookDecodeFunc{func(in reflect.Type, out reflect.Type, data any) (any, error) {
			if out == reflect.TypeFor[time.Duration]() {
				switch in.Kind() {
				case reflect.String:
					return time.ParseDuration(data.(string))
				}
			}

			return data, nil
		}},
		WeaklyTypedInput:      true,
		WeaklyIgnoreSeperator: true,
		WeaklyDashUnderscore:  true,
	}

	cfgMap := map[string]any{
		"name":  "Altay",
		"count": "42",
		"sub": map[string]any{
			"enabled": "true",
		},
		"time": "1h30m",
	}

	var cfg Config

	if err := decoder.Decode(cfgMap, &cfg); err != nil {
		fmt.Printf("decode error: %v", err)
		return
	}

	fmt.Printf("%v %v %v %v", cfg.Name, cfg.Count, cfg.Sub.Enabled, cfg.Time)
	// Output:
	// Altay 42 true 1h30m0s
}
