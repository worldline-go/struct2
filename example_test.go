package struct2_test

import (
	"fmt"
	"sort"
	"time"

	"github.com/worldline-go/struct2"
	"github.com/worldline-go/struct2/types"
)

func SortPrint(m map[string]interface{}) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("Type: %T, Value: %v\n", m[k], m[k])
	}
}

func Example() {
	type ColorGroup struct {
		ID     int        `json:"id"`
		Name   string     `json:"name"`
		Colors []string   `json:"colors"`
		Date   types.Time `json:"time"`
	}

	d, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")

	group := ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
		Date:   types.Time{Time: d},
	}

	result := new(struct2.Decoder).SetTagName("json").Map(group) // default tag name is "struct"

	// fmt.Printf("%#v", result)
	SortPrint(result)
	// Output:
	// Type: []string, Value: [Crimson Red Ruby Maroon]
	// Type: int, Value: 1
	// Type: string, Value: Reds
	// Type: time.Time, Value: 2006-01-02 15:04:05 +0000 UTC
}
