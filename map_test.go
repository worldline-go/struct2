package struct2

import (
	"reflect"
	"testing"
)

type stringerValue int

func (s stringerValue) String() string {
	return "stringerValue"
}

func str2Ptr(s string) *string {
	return &s
}

func int2Ptr(i int) *int {
	return &i
}

func TestDecoder_Map(t *testing.T) {
	type Train struct {
		Wagon *int `struct:"wagon,ptr2"`
	}
	type TrainNoPtr2 struct {
		Wagon *int `struct:"wagon"`
	}

	type TrainStringTags struct {
		Wagon    *int          `struct:"wagon,string"`
		Stringer stringerValue `struct:"stringer,string"`
	}

	train := Train{
		Wagon: int2Ptr(5),
	}

	type args struct {
		s interface{}
	}
	tests := []struct {
		name    string
		decoder Decoder
		args    args
		want    map[string]interface{}
	}{

		{
			name: "unexported test",
			args: args{
				s: struct {
					Name       string `struct:"name"`
					unexported string
				}{
					Name:       "abc",
					unexported: "unexported",
				},
			},
			want: map[string]interface{}{
				"name": "abc",
			},
		},
		{
			name: "simple test",
			args: args{
				s: struct {
					Name string `struct:"name"`
					Ptr  *string
				}{
					Name: "abc",
					Ptr:  str2Ptr("pointer"),
				},
			},
			want: map[string]interface{}{
				"name": "abc",
				"Ptr":  str2Ptr("pointer"),
			},
		},
		{
			name: "nil test",
			args: args{
				s: (*struct {
					Name string  `struct:"name"`
					Ptr  *string `struct:"ptr,ptr2"`
				})(nil),
			},
			want: nil,
		},
		{
			name: "ptr2 test",
			args: args{
				s: struct {
					Name string  `struct:"name"`
					Ptr  *string `struct:"ptr,ptr2"`
				}{
					Name: "abc",
					Ptr:  str2Ptr("pointer"),
				},
			},
			want: map[string]interface{}{
				"name": "abc",
				"ptr":  "pointer",
			},
		},
		{
			name: "deep test",
			args: args{
				s: struct {
					Name  string  `struct:"name"`
					Ptr   *string `struct:"ptr,ptr2"`
					Train struct {
						Wagon *int `struct:"wagon,ptr2"`
					} `struct:"train"`
				}{
					Name: "abc",
					Ptr:  str2Ptr("pointer"),
					Train: struct {
						Wagon *int `struct:"wagon,ptr2"`
					}{Wagon: int2Ptr(5)},
				},
			},
			want: map[string]interface{}{
				"name": "abc",
				"ptr":  "pointer",
				"train": map[string]interface{}{
					"wagon": 5,
				},
			},
		},
		{
			name: "deep test with omitnested",
			args: args{
				s: struct {
					Name  string  `struct:"name"`
					Ptr   *string `struct:"ptr,ptr2"`
					Train *Train  `struct:"train,omitnested,ptr2"`
				}{
					Name:  "abc",
					Ptr:   str2Ptr("pointer"),
					Train: &train,
				},
			},
			want: map[string]interface{}{
				"name":  "abc",
				"ptr":   "pointer",
				"train": train,
			},
		},
		{
			name: "slice pointer of struct",
			args: args{
				s: struct {
					Name  string   `struct:"name"`
					Ptr   *string  `struct:"ptr,ptr2"`
					Train *[]Train `struct:"train"`
				}{
					Name:  "abc",
					Ptr:   str2Ptr("pointer"),
					Train: &[]Train{train},
				},
			},
			want: map[string]interface{}{
				"name":  "abc",
				"ptr":   "pointer",
				"train": []interface{}{map[string]interface{}{"wagon": 5}},
			},
		},
		{
			name:    "pointer with ForcePtr2",
			decoder: Decoder{ForcePtr2: true},
			args: args{
				s: struct {
					Name  string      `struct:"name"`
					Ptr   *string     `struct:"ptr"`
					Train TrainNoPtr2 `struct:"train"`
				}{
					Name:  "abc",
					Ptr:   str2Ptr("pointer"),
					Train: TrainNoPtr2{},
				},
			},
			want: map[string]interface{}{
				"name": "abc",
				"ptr":  "pointer",
				"train": map[string]interface{}{
					"wagon": 0,
				},
			},
		},
		{
			name:    "pointer with OmitNilPtr",
			decoder: Decoder{OmitNilPtr: true},
			args: args{
				s: struct {
					Name     string   `struct:"name"`
					Ptr      *string  `struct:"ptr,ptr2"`
					Train    Train    `struct:"train"`
					TrainPtr *Train   `struct:"trainPtr"`
					Trains   *[]Train `struct:"trains"`
				}{
					Name:     "abc",
					Ptr:      nil,
					Train:    Train{},
					TrainPtr: &Train{},
					Trains: &[]Train{
						{},
					},
				},
			},
			want: map[string]interface{}{
				"name":     "abc",
				"train":    map[string]interface{}{},
				"trainPtr": map[string]interface{}{},
				"trains": []interface{}{
					map[string]interface{}{},
				},
			},
		},
		{
			name: "ouput with OuputCamelCase",
			decoder: Decoder{
				OuputCamelCase: true,
			},
			args: args{
				s: struct {
					Name     string
					Name2    string `struct:"Name2"`
					Train    Train
					TrainPtr *Train
					Trains   []Train
				}{
					Name:  "abc",
					Name2: "def",
					Train: Train{
						Wagon: int2Ptr(5),
					},
					TrainPtr: &Train{
						Wagon: int2Ptr(5),
					},
					Trains: []Train{
						{
							Wagon: int2Ptr(5),
						},
					},
				},
			},
			want: map[string]interface{}{
				"name":  "abc",
				"Name2": "def",
				"train": map[string]interface{}{
					"wagon": 5,
				},
				"trainPtr": map[string]interface{}{
					"wagon": 5,
				},
				"trains": []interface{}{
					map[string]interface{}{
						"wagon": 5,
					},
				},
			},
		},
		{
			name: "string tags",
			args: args{
				s: struct {
					Train TrainStringTags `struct:"train"`
				}{
					Train: TrainStringTags{
						Wagon:    int2Ptr(5),
						Stringer: 1,
					},
				},
			},
			want: map[string]interface{}{
				"train": map[string]interface{}{
					"wagon":    "5",
					"stringer": "stringerValue",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.decoder.Map(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decoder.Map() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestDecoder_MapOmitNested(t *testing.T) {
	type Train struct {
		Wagon *int `struct:"wagon,ptr2"`
	}

	train := Train{
		Wagon: int2Ptr(5),
	}

	type args struct {
		s interface{}
	}
	tests := []struct {
		name    string
		decoder Decoder
		args    args
		want    map[string]interface{}
	}{
		{
			name: "deep test with omitnested",
			args: args{
				s: struct {
					Name  string  `struct:"name"`
					Ptr   *string `struct:"ptr,ptr2"`
					Train *Train  `struct:"train,ptr2"`
				}{
					Name:  "abc",
					Ptr:   str2Ptr("pointer"),
					Train: &train,
				},
			},
			want: map[string]interface{}{
				"name":  "abc",
				"ptr":   "pointer",
				"train": train,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.decoder.MapOmitNested(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decoder.Map() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestDecoder_CustomHook(t *testing.T) {
	type args struct {
		s interface{}
	}
	tests := []struct {
		name    string
		decoder Decoder
		args    args
		want    map[string]interface{}
	}{
		{
			name: "simple test",
			decoder: Decoder{
				Hooks: []HookFunc{func(v reflect.Value) (interface{}, error) {
					if v.Kind() == reflect.String {
						return "str_" + v.Interface().(string), nil
					}
					return nil, ErrContinueHook
				}},
			},
			args: args{
				s: struct {
					Name string `struct:"name"`
				}{
					Name: "abc",
				},
			},
			want: map[string]interface{}{
				"name": "str_abc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.decoder.Map(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decoder.Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
