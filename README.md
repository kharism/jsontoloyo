# JSON Mapping-Traversing library

this library can traverse and do simple mapping of 1 json (map[string]interface{}) to another json. For now it can't deals with aggregating json array into single json.

Example for simple traversing

```
    input := "$.name.first"
	source := strings.NewReader(input)
	parser := jsontoloyo.NewParser(source)
	selectors := []*jsontoloyo.Selector{}
	sourceMap := map[string]interface{}{}
	sourceMap["name"] = map[string]interface{}{"first": "LL", "last": "XX"}

	err := parser.Parse(&selectors)
	if err != nil {
		os.Exit(-1)
	}
	fmt.Println(selectors)
	if len(selectors) != 3 {
		os.Exit(-1)
	}
	traverser := jsontoloyo.NewTraverser(sourceMap, selectors)
	kk, err := traverser.Traverse()
	if err != nil {
		os.Exit(-1)
	}
	fmt.Println(kk) // this should output LL
```

Get specific index of array

```
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
		fmt.Println(err.Error())
		os.Exit(-1)
	}
    fmt.Println(kk)// kk should be 1
```
get aggregate of an array
```
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
    fmt.Println(kk)//kk should be 10
```
basic mapping json to json

```
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
```