package struct2

import (
	"reflect"
	"testing"
)

func Test_parseTag(t *testing.T) {
	type args struct {
		tag string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 tagOptions
	}{
		{
			name: "empty",
			args: args{
				tag: "",
			},
			want:  "",
			want1: tagOptions{},
		},
		{
			name: "multi",
			args: args{
				tag: "struct,omitempty,ptr2",
			},
			want:  "struct",
			want1: tagOptions{"omitempty", "ptr2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := parseTag(tt.args.tag)
			if got != tt.want {
				t.Errorf("parseTag() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseTag() got1 = %v, want1 %v", got1, tt.want1)
			}
		})
	}
}
