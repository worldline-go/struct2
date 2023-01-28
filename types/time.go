package types

import (
	"bytes"
	"fmt"
	"time"

	"github.com/worldline-go/struct2"
)

type Time struct {
	time.Time
}

var _ struct2.Hooker = Time{}

// UnmarshalJSON custom for fit other time layouts.
func (t *Time) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		return nil
	}

	for _, layout := range []string{time.RFC3339, time.RFC3339Nano, "2006-01-02 15:04:05.000000"} {
		if parsedTime, err := time.Parse(`"`+layout+`"`, string(b)); err == nil {
			*t = Time{parsedTime}

			return nil
		}
	}

	return fmt.Errorf("cannot unmarshal; unknown time layout %s", b)
}

// Struct2Hook hook function for decode struct2.
func (t Time) Struct2Hook() interface{} {
	return t.Time
}
