package jsontoloyo_test

import (
	"fmt"
	"github/kharism/jsontoloyo"
	"os"
	"strings"
	"testing"
)

func TestParseBasicField(t *testing.T) {
	input := "$.name.first"
	source := strings.NewReader(input)
	parser := jsontoloyo.NewParser(source)
	selectors := []*jsontoloyo.Selector{}
	err := parser.Parse(&selectors)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	fmt.Println(selectors)
	if len(selectors) != 3 {
		t.Fail()
	}
}
func TestParseBasicArray(t *testing.T) {
	input := "$.ii[1]"
	source := strings.NewReader(input)
	parser := jsontoloyo.NewParser(source)
	selectors := []*jsontoloyo.Selector{}
	err := parser.Parse(&selectors)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	fmt.Println(selectors)
	if len(selectors) != 5 {
		t.Fail()
	}
}
func TestParseComplexIndexing(t *testing.T) {
	input := "$.ii[1].name"
	source := strings.NewReader(input)
	parser := jsontoloyo.NewParser(source)
	// sourceMap["ii"] = []interface{}{0, map[string]interface{}{"name": "KK"}, 2, 3, 4}
	selectors := []*jsontoloyo.Selector{}
	err := parser.Parse(&selectors)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	fmt.Println(selectors)
	if len(selectors) != 6 {
		t.Fail()
	}
}
func TestParseIterable(t *testing.T) {
	input := "$[]"
	source := strings.NewReader(input)
	parser := jsontoloyo.NewParser(source)
	selectors := []*jsontoloyo.Selector{}
	err := parser.Parse(&selectors)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	fmt.Println(selectors)
	if len(selectors) != 3 {
		t.Fail()
	}

	input = "$[10]"
	source = strings.NewReader(input)
	parser = jsontoloyo.NewParser(source)
	selectors = []*jsontoloyo.Selector{}
	err = parser.Parse(&selectors)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	fmt.Println(selectors)
	if len(selectors) != 4 {
		t.Fail()
	}
}

func TestParseAggregate(t *testing.T) {
	input := "$.$avg($.age)"
	source := strings.NewReader(input)
	parser := jsontoloyo.NewParser(source)
	selectors := []*jsontoloyo.Selector{}
	err := parser.Parse(&selectors)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	fmt.Println(selectors)
	if len(selectors) != 2 {
		t.Fail()
	}
}
