package struct2

import (
	"reflect"
	"testing"
)

func TestPtr2Concrete(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "non pointer test",
			args: args{
				val: "non pointer",
			},
			want: "non pointer",
		},
		{
			name: "pointer test",
			args: args{
				val: str2Ptr("pointer"),
			},
			want: "pointer",
		},
		{
			name: "struct test",
			args: args{
				val: struct{ name string }{name: "struct"},
			},
			want: struct{ name string }{name: "struct"},
		},
		{
			name: "pointer struct test",
			args: args{
				val: &struct{ name string }{name: "struct"},
			},
			want: struct{ name string }{name: "struct"},
		},
		{
			name: "pointer struct test with nil",
			args: args{
				val: (*struct{ name string })(nil),
			},
			want: struct{ name string }{},
		},
		{
			name: "nil string",
			args: args{
				val: (*string)(nil),
			},
			want: "",
		},
		{
			name: "nil int",
			args: args{
				val: (*int)(nil),
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ptr2Concrete(tt.args.val); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ptr2Concrete() = %v, want %v", got, tt.want)
			}
		})
	}
}
