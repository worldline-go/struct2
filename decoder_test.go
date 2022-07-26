package struct2

import (
	"reflect"
	"testing"
)

func Test_interface2StructValue(t *testing.T) {
	type args struct {
		s interface{}
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name: "panic not struct value",
			args: args{
				s: "not struct",
			},
			wantPanic: true,
		},
		{
			name: "struct value",
			args: args{
				s: struct {
					name int
				}{
					name: 123,
				},
			},
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.wantPanic {
					t.Errorf("panic %t, wantPanic %t", (r != nil), tt.wantPanic)
				}
			}()

			value2StructValue(reflect.ValueOf(tt.args.s))
		})
	}
}
