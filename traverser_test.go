package jsontoloyo_test

import (
	"fmt"
	"github/kharism/jsontoloyo"
	"os"
	"strings"
	"testing"
)

func TestTraverseBasicField(t *testing.T) {
	input := "$.name.first"
	source := strings.NewReader(input)
	parser := jsontoloyo.NewParser(source)
	selectors := []*jsontoloyo.Selector{}
	sourceMap := map[string]interface{}{}
	sourceMap["name"] = map[string]interface{}{"first": "LL", "last": "XX"}

	err := parser.Parse(&selectors)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	fmt.Println(selectors)
	if len(selectors) != 3 {
		t.Log("selector len is", len(selectors))
		t.Fail()
	}
	traverser := jsontoloyo.NewTraverser(sourceMap, selectors)
	kk, err := traverser.Traverse()
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	if kk != "LL" {
		t.Log("Content is", kk)
		t.Fail()
	}
	sourceMap["name"] = map[string]interface{}{"first": map[string]interface{}{"kk": "LL"}, "last": "XX"}
	traverser = jsontoloyo.NewTraverser(sourceMap, selectors)
	kk, err = traverser.Traverse()
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	if kk.(map[string]interface{})["kk"] != "LL" {
		t.Fail()
	}
}

func TestTraverseArray(t *testing.T) {
	input := "$.ii[1]"
	source := strings.NewReader(input)
	parser := jsontoloyo.NewParser(source)
	selectors := []*jsontoloyo.Selector{}
	sourceMap := map[string]interface{}{}
	sourceMap["ii"] = []int{0, 1, 2, 3, 4}
	err := parser.Parse(&selectors)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	traverser := jsontoloyo.NewTraverser(sourceMap, selectors)
	kk, err := traverser.Traverse()
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	if kk != 1 {
		t.Log("Content is", kk)
		t.Fail()
	}

	/// case 2
	input = "$.ii[1].name"
	source = strings.NewReader(input)
	parser = jsontoloyo.NewParser(source)
	sourceMap["ii"] = []interface{}{0, map[string]interface{}{"name": "KK"}, 2, 3, 4}
	selectors = []*jsontoloyo.Selector{}
	err = parser.Parse(&selectors)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	traverser = jsontoloyo.NewTraverser(sourceMap, selectors)
	kk, err = traverser.Traverse()
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	if kk != "KK" {
		t.Log("Content is", kk)
		t.Fail()
	}

	/// case 3
	input = "$.ii[1].name[0]"
	source = strings.NewReader(input)
	parser = jsontoloyo.NewParser(source)
	sourceMap["ii"] = []interface{}{0, map[string]interface{}{"name": []float64{0.1, 0.2}}, 2, 3, 4}
	selectors = []*jsontoloyo.Selector{}
	err = parser.Parse(&selectors)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	traverser = jsontoloyo.NewTraverser(sourceMap, selectors)
	kk, err = traverser.Traverse()
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	if kk != 0.1 {
		t.Log("Content is", kk)
		t.Fail()
	}

	//case 4
	input = "$.ii[1].name[0].aa"
	source = strings.NewReader(input)
	parser = jsontoloyo.NewParser(source)
	sourceMap["ii"] = []interface{}{0, map[string]interface{}{"name": []interface{}{map[string]interface{}{"aa": "ii"}, 0.2}}, 2, 3, 4}
	selectors = []*jsontoloyo.Selector{}
	err = parser.Parse(&selectors)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	traverser = jsontoloyo.NewTraverser(sourceMap, selectors)
	kk, err = traverser.Traverse()
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	if kk != "ii" {
		t.Log("Content is", kk)
		t.Fail()
	}

}

func TestTraverseAggregate(t *testing.T) {
	input := "$.$sum($.ii)"
	source := strings.NewReader(input)
	parser := jsontoloyo.NewParser(source)
	selectors := []*jsontoloyo.Selector{}
	sourceMap := map[string]interface{}{}
	sourceMap["ii"] = []int{0, 1, 2, 3, 4}
	err := parser.Parse(&selectors)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	if len(selectors) != 2 {
		t.Fail()
	}
	traverser := jsontoloyo.NewTraverser(sourceMap, selectors)
	kk, err := traverser.Traverse()
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	if kk != 10.0 {
		t.Log("Content is", kk)
		t.Fail()
	}

	// case 2
	input = "$.ii[1].$sum($.name)"
	source = strings.NewReader(input)
	parser = jsontoloyo.NewParser(source)
	sourceMap["ii"] = []interface{}{0, map[string]interface{}{"name": []float64{0.1, 0.2}}, 2, 3, 4}
	selectors = []*jsontoloyo.Selector{}
	err = parser.Parse(&selectors)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	traverser = jsontoloyo.NewTraverser(sourceMap, selectors)
	kk, err = traverser.Traverse()
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	if kk != 0.30000000000000004 {
		t.Log("Content is", kk)
		t.Fail()
	}
}
