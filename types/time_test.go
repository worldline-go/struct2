package types

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/worldline-go/struct2"
)

func TestTime_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Date Time `json:"date"`
	}

	tests := []struct {
		name    string
		fields  fields
		arg     string
		wantErr bool
		want    string
	}{
		{
			name:    "RFC3339 test",
			fields:  fields{},
			arg:     "2021-10-15T16:59:56Z",
			wantErr: false,
			want:    `{"date":"2021-10-15T16:59:56Z"}`,
		},
		{
			name:    "layout 2006-01-02 15:04:05.000000",
			fields:  fields{},
			arg:     "2022-02-22 22:22:22.000000",
			wantErr: false,
			want:    `{"date":"2022-02-22T22:22:22Z"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := json.Unmarshal([]byte(`{"date":"`+tt.arg+`"}`), &tt.fields); (err != nil) != tt.wantErr {
				t.Errorf("Time.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				v, _ := json.Marshal(tt.fields)
				if string(v) != tt.want {
					t.Errorf("result = %s, want = %s", v, tt.arg)
				}
			}
		})
	}
}

func TestDecoder_Map(t *testing.T) {
	type args struct {
		s interface{}
	}
	tests := []struct {
		name    string
		decoder struct2.Decoder
		args    args
		want    map[string]interface{}
	}{
		{
			name: "simple test",
			decoder: struct2.Decoder{
				TagName: "db",
			},
			args: args{
				s: struct {
					Date  Time      `db:"date"`
					Date2 time.Time `db:"date2"`
					Date3 *Time     `db:"date3,ptr2"`
					Date4 *Time     `db:"date4"`
				}{
					Date: func() Time {
						d, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
						return Time{Time: d}
					}(),
					Date2: func() time.Time {
						d, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
						return d
					}(),
					Date3: func() *Time {
						d, _ := time.Parse(time.RFC3339, "2016-01-02T15:04:05Z")
						return &Time{Time: d}
					}(),
					Date4: func() *Time {
						d, _ := time.Parse(time.RFC3339, "2016-01-02T15:04:05Z")
						return &Time{Time: d}
					}(),
				},
			},
			want: map[string]interface{}{
				"date": func() time.Time {
					d, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
					return d
				}(),
				"date2": func() time.Time {
					d, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
					return d
				}(),
				"date3": func() time.Time {
					d, _ := time.Parse(time.RFC3339, "2016-01-02T15:04:05Z")
					return d
				}(),
				"date4": func() time.Time {
					d, _ := time.Parse(time.RFC3339, "2016-01-02T15:04:05Z")
					return d
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.decoder.Map(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				// fmt.Printf("got: %T, want: %T\n", got["date3"], tt.want["date3"])
				t.Errorf("Decoder.Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
