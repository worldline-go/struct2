package struct2

import (
	"reflect"
	"testing"
)

func TestDecoder_GetFields(t *testing.T) {
	type args struct {
		s interface{}
	}
	tests := []struct {
		name    string
		decoder Decoder
		args    args
		want    []string
	}{
		{
			name: "simple sturct",
			args: args{
				s: struct {
					test  string
					Test2 int
				}{},
			},
			want: []string{"Test2"},
		},
		{
			name: "unset all",
			args: args{
				s: struct {
					test  string `struct:"testing"`
					Test2 int    `struct:"-"`
				}{},
			},
			want: []string{},
		},
		{
			name: "with tag name",
			args: args{
				s: struct {
					test  string
					Test2 int `struct:"test_2"`
				}{},
			},
			want: []string{"test_2"},
		},
		{
			name: "nil",
			args: args{
				s: (*struct {
					test  string
					Test2 int `struct:"test_2"`
				})(nil),
			},
			want: []string{"test_2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if v := tt.decoder.GetFields(tt.args.s); !reflect.DeepEqual(v, tt.want) {
				t.Errorf("GetFields want %v, got %v", tt.want, v)
			}
		})
	}
}
