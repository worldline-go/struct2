package struct2_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/worldline-go/struct2"
	"github.com/worldline-go/struct2/types"
)

type ColorGroup struct {
	ID     int      `db:"id"`
	Name   string   `db:"name"`
	Colors []string `db:"colors"`
	// custom type with implemented Hooker interface
	// covertion result to time.Time
	Date types.Time `db:"time"`
	// CustomValue will be pointer to int
	CustomValue *int `db:"custom_value,ptr2"`
	// Inner field will be flatten, so override ID from root
	Inner Inner `db:"inner,flatten"`
	// Coordinate field will be same struct, not touching
	Coordinate Coordinate `db:"coordinate,omitnested"`
	// Secret field will be ignored
	Secret string `cfg:"-"`
}

type Inner struct {
	ID int `db:"id"`
}

type Coordinate struct {
	X int `db:"x"`
	Y int `db:"y"`
}

type ColorGroupData struct {
	Time time.Time
}

func (d ColorGroupData) Data() interface{} {
	v := 5
	return ColorGroup{
		ID:          1,
		Name:        "Reds",
		Colors:      []string{"Crimson", "Red", "Ruby", "Maroon"},
		Date:        types.Time{Time: d.Time},
		CustomValue: &v,
		Inner: Inner{
			ID: 2,
		},
		Coordinate: Coordinate{
			X: 1,
			Y: 2,
		},
		Secret: "secret",
	}
}

func (d ColorGroupData) MapData() interface{} {
	return map[string]interface{}{
		"id":           2,
		"name":         "Reds",
		"colors":       []string{"Crimson", "Red", "Ruby", "Maroon"},
		"time":         d.Time,
		"custom_value": 5,
		"coordinate": Coordinate{
			X: 1,
			Y: 2,
		},
	}
}

// defination of example data

type exampleData interface {
	Data() interface{}
	MapData() interface{}
}

var mapExamples = map[string]exampleData{
	"ColorGroup": ColorGroupData{
		Time: time.Now().Truncate(time.Second),
	},
}

// run specific DATA by name
func Test_Data(t *testing.T) {
	decoder := struct2.Decoder{
		TagName:               "db",
		BackupTagName:         "cfg",
		WeaklyDashUnderscore:  true,
		WeaklyIgnoreSeperator: true,
		WeaklyTypedInput:      true,
	}

	example := os.Getenv("DATA")
	if example == "" {
		t.Skip("DATA env variable is empty")
	}

	mapper, ok := mapExamples[example]
	if !ok {
		t.Fatalf("example %s not found", example)
	}

	inputData := mapper.Data()
	outputData := decoder.Map(inputData)
	correctData := mapper.MapData()

	fmt.Printf("\n# input\n%#v\n", inputData)
	fmt.Printf("\n# output\n%#v\n", outputData)
	fmt.Printf("\n# correction\n%#v\n", correctData)

	if !reflect.DeepEqual(outputData, correctData) {
		t.Fatalf("output data not equal to correct data %v != %v", outputData, correctData)
	}
}
