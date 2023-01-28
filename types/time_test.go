package types

import (
	"bytes"
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
		arg     []byte
		wantErr bool
		want    []byte
	}{
		{
			name:    "RFC3339 test",
			fields:  fields{},
			arg:     []byte(`{"date":"2021-10-15T16:59:56Z"}`),
			wantErr: false,
			want:    []byte(`{"date":"2021-10-15T16:59:56Z"}`),
		},
		{
			name:    "layout 2006-01-02 15:04:05.000000",
			fields:  fields{},
			arg:     []byte(`{"date":"2022-02-22T22:22:22Z"}`),
			wantErr: false,
			want:    []byte(`{"date":"2022-02-22T22:22:22Z"}`),
		},
		{
			name:    "null test",
			fields:  fields{},
			arg:     []byte("null"),
			wantErr: false,
			want:    []byte(`{"date":"0001-01-01T00:00:00Z"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := json.Unmarshal(tt.arg, &tt.fields); (err != nil) != tt.wantErr {
				t.Errorf("Time.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				v, _ := json.Marshal(tt.fields)
				if !bytes.Equal(v, tt.want) {
					t.Errorf("result = %s, want = %s", v, tt.arg)
				}
			}
		})
	}
}

func TestDecoder_Map(t *testing.T) {
	tests := []struct {
		name    string
		decoder struct2.Decoder
		arg     interface{}
		want    map[string]interface{}
	}{
		{
			name: "simple test",
			decoder: struct2.Decoder{
				TagName: "db",
			},
			arg: struct {
				Date    Time      `db:"date"`
				Date2   time.Time `db:"date2"`
				Date3   *Time     `db:"date3,ptr2"`
				Date4   *Time     `db:"date4"`
				Date5   *Time     `db:"date5"`
				DateStr *string   `db:"DateStr,ptr2"`
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
				Date5: nil,
				DateStr: func() *string {
					v := "2016-01-02T15:04:05Z"
					return &v
				}(),
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
				"date5":   nil,
				"DateStr": "2016-01-02T15:04:05Z",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.decoder.Map(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decoder.Map() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
