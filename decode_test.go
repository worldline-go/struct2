package struct2

import (
	"reflect"
	"testing"
)

func TestDecoder_Decode(t *testing.T) {
	type fields struct {
		TagName               string
		Hooks                 []HookFunc
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
		wantErr bool
		want    interface{}
	}{
		{
			name:   "test",
			fields: fields{},
			args: args{
				input: map[string]interface{}{
					"test": "test",
				},
				output: &struct {
					Test string `struct:"test"`
				}{},
			},
			wantErr: false,
			want: &struct {
				Test string `struct:"test"`
			}{
				Test: "test",
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
			wantErr: false,
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
			wantErr: false,
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
			wantErr: false,
			want: &struct {
				Test string `struct:"test X"`
			}{
				Test: "test",
			},
		},
		{
			name: "copy []byte",
			args: args{
				input:  append(make([]byte, 0, 100), []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}...),
				output: func() interface{} { v := make([]byte, 0, 10); return &v }(),
			},
			wantErr: false,
			want:    func() interface{} { v := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}; return &v }(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Decoder{
				TagName:               tt.fields.TagName,
				Hooks:                 tt.fields.Hooks,
				WeaklyTypedInput:      tt.fields.WeaklyTypedInput,
				ZeroFields:            tt.fields.ZeroFields,
				Squash:                tt.fields.Squash,
				IgnoreUntaggedFields:  tt.fields.IgnoreUntaggedFields,
				BackupTagName:         tt.fields.BackupTagName,
				WeaklyDashUnderscore:  tt.fields.WeaklyDashUnderscore,
				WeaklyIgnoreSeperator: tt.fields.WeaklyIgnoreSeperator,
			}
			if err := d.Decode(tt.args.input, tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.Decode() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(tt.args.output, tt.want) {
				t.Errorf("Decoder.Decode() = %v, want %v", tt.args.output, tt.want)
			}
		})
	}
}
