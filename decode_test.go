package struct2

import (
	"reflect"
	"testing"
	"time"
)

func TestDecoder_Decode(t *testing.T) {
	timeNow := time.Now()

	type fields struct {
		TagName               string
		Hooks                 []HookFunc
		HooksDecode           []HookDecodeFunc
		WeaklyTypedInput      bool
		ZeroFields            bool
		Squash                bool
		IgnoreUntaggedFields  bool
		BackupTagName         string
		WeaklyDashUnderscore  bool
		WeaklyIgnoreSeperator bool
	}
	type args struct {
		input  interface{}
		output interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
		want    interface{}
	}{
		{
			name: "test",
			args: args{
				input: map[string]interface{}{
					"test": "test",
				},
				output: &struct {
					Test string `struct:"test"`
				}{},
			},
			wantErr: nil,
			want: &struct {
				Test string `struct:"test"`
			}{
				Test: "test",
			},
		},
		{
			name: "test",
			fields: fields{
				WeaklyTypedInput: true,
			},
			args: args{
				input: map[string]interface{}{
					"string":  "test",
					"bool":    "True",
					"int":     1,
					"float64": 1.1,
					"float32": func() interface{} { return 1.1 }(),
					"byte":    1,
					"rune":    '#',
					"uint":    1,
					"uint8":   1,
					"uint16":  1,
					"uint32":  1,
					"uint64":  1,
					"int8":    1,
					"int16":   1,
					"int32":   1,
					"int64":   1,
					"time":    func() interface{} { return timeNow }(),
				},
				output: &struct {
					String  string    `struct:"string"`
					Bool    bool      `struct:"bool"`
					Int     int       `struct:"int"`
					Float64 float64   `struct:"float64"`
					Float32 float32   `struct:"float32"`
					Byte    byte      `struct:"byte"`
					Rune    rune      `struct:"rune"`
					Uint    uint      `struct:"uint"`
					Uint8   uint8     `struct:"uint8"`
					Uint16  uint16    `struct:"uint16"`
					Uint32  uint32    `struct:"uint32"`
					Uint64  uint64    `struct:"uint64"`
					Int8    int8      `struct:"int8"`
					Int16   int16     `struct:"int16"`
					Int32   int32     `struct:"int32"`
					Int64   int64     `struct:"int64"`
					Time    time.Time `struct:"time"`
				}{},
			},
			wantErr: nil,
			want: &struct {
				String  string    `struct:"string"`
				Bool    bool      `struct:"bool"`
				Int     int       `struct:"int"`
				Float64 float64   `struct:"float64"`
				Float32 float32   `struct:"float32"`
				Byte    byte      `struct:"byte"`
				Rune    rune      `struct:"rune"`
				Uint    uint      `struct:"uint"`
				Uint8   uint8     `struct:"uint8"`
				Uint16  uint16    `struct:"uint16"`
				Uint32  uint32    `struct:"uint32"`
				Uint64  uint64    `struct:"uint64"`
				Int8    int8      `struct:"int8"`
				Int16   int16     `struct:"int16"`
				Int32   int32     `struct:"int32"`
				Int64   int64     `struct:"int64"`
				Time    time.Time `struct:"time"`
			}{
				String:  "test",
				Bool:    true,
				Int:     1,
				Float64: 1.1,
				Float32: 1.1,
				Byte:    1,
				Rune:    '#',
				Uint:    1,
				Uint8:   1,
				Uint16:  1,
				Uint32:  1,
				Uint64:  1,
				Int8:    1,
				Int16:   1,
				Int32:   1,
				Int64:   1,
				Time:    timeNow,
			},
		},
		{
			name: "nil struct",
			args: args{
				input: map[string]interface{}{
					"abc": "x",
				},
				output: &struct {
					Test *struct {
						Abc string `struct:"abc"`
					} `struct:"test"`
				}{},
			},
			wantErr: nil,
			want: &struct {
				Test *struct {
					Abc string `struct:"abc"`
				} `struct:"test"`
			}{
				Test: nil,
			},
		},
		{
			name: "hook test",
			fields: fields{
				Hooks: []HookFunc{
					func(v reflect.Value) (interface{}, error) {
						if v.Kind() == reflect.String {
							return v.Interface().(string) + "_hooked", nil
						}

						return nil, ErrContinueHook
					},
				},
			},
			args: args{
				input: map[string]interface{}{
					"test": "test",
				},
				output: &struct {
					Test string `struct:"test"`
				}{},
			},
			wantErr: nil,
			want: &struct {
				Test string `struct:"test"`
			}{
				Test: "test_hooked",
			},
		},
		{
			name: "WeaklyDashUnderscore",
			fields: fields{
				WeaklyDashUnderscore: true,
			},
			args: args{
				input: map[string]interface{}{
					"test_x": "test",
				},
				output: &struct {
					Test string `struct:"test-x"`
				}{},
			},
			wantErr: nil,
			want: &struct {
				Test string `struct:"test-x"`
			}{
				Test: "test",
			},
		},
		{
			name: "WeaklyIgnoreSeperator",
			fields: fields{
				WeaklyIgnoreSeperator: true,
			},
			args: args{
				input: map[string]interface{}{
					"test_x": "test",
				},
				output: &struct {
					Test string `struct:"test X"`
				}{},
			},
			wantErr: nil,
			want: &struct {
				Test string `struct:"test X"`
			}{
				Test: "test",
			},
		},
		{
			name: "hooksdecode",
			fields: fields{
				WeaklyIgnoreSeperator: true,
				HooksDecode: []HookDecodeFunc{
					func(t1, t2 reflect.Type, data interface{}) (interface{}, error) {
						if t2 != reflect.TypeOf(time.Duration(0)) {
							return data, nil
						}

						switch t1.Kind() {
						case reflect.String:
							return time.ParseDuration(data.(string))
						case reflect.Int:
							return time.Duration(data.(int)), nil
						case reflect.Int64:
							return time.Duration(data.(int64)), nil
						case reflect.Float64:
							return time.Duration(data.(float64)), nil
						default:
							return data, nil
						}
					},
				},
			},
			args: args{
				input: map[string]interface{}{
					"test_x": "5s",
					"test_y": 1_000_000_000,
				},
				output: &struct {
					Test  time.Duration `struct:"test X"`
					TestY time.Duration `struct:"test_y"`
				}{},
			},
			wantErr: nil,
			want: &struct {
				Test  time.Duration `struct:"test X"`
				TestY time.Duration `struct:"test_y"`
			}{
				Test:  time.Duration(5 * time.Second),
				TestY: time.Duration(1 * time.Second),
			},
		},
		{
			name: "copy []byte",
			args: args{
				input:  append(make([]byte, 0, 100), []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}...),
				output: func() interface{} { v := make([]byte, 0, 10); return &v }(),
			},
			wantErr: nil,
			want:    func() interface{} { v := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}; return &v }(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Decoder{
				TagName:               tt.fields.TagName,
				Hooks:                 tt.fields.Hooks,
				HooksDecode:           tt.fields.HooksDecode,
				WeaklyTypedInput:      tt.fields.WeaklyTypedInput,
				ZeroFields:            tt.fields.ZeroFields,
				Squash:                tt.fields.Squash,
				IgnoreUntaggedFields:  tt.fields.IgnoreUntaggedFields,
				BackupTagName:         tt.fields.BackupTagName,
				WeaklyDashUnderscore:  tt.fields.WeaklyDashUnderscore,
				WeaklyIgnoreSeperator: tt.fields.WeaklyIgnoreSeperator,
			}
			if err := d.Decode(tt.args.input, tt.args.output); !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Decoder.Decode() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(tt.args.output, tt.want) {
				t.Errorf("Decoder.Decode() = %v, want %v", tt.args.output, tt.want)
			}
		})
	}
}
