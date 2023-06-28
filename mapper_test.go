package jsontoloyo_test

import (
	"github/kharism/jsontoloyo"
	"testing"
)

func TestMapping(t *testing.T) {
	source := map[string]interface{}{
		"scores": []int{10, 10, 40},
	}
	sink := map[string]interface{}{}

	mappings := jsontoloyo.NewMappers()
	mappings.AddMapping("avg_score", "$.$avg($.scores)")
	err := mappings.Mapping(source, sink)
	if err != nil {
		t.Log("err", err.Error())
		t.Fail()
	}
	if _, ok := sink["avg_score"]; !ok {
		t.Log("No new field")
		t.FailNow()
	}
	if sink["avg_score"] != 20.0 {
		t.Log("avg", sink["avg_score"])
		t.Fail()
	}
}
