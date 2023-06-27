package jsontoloyo

import (
	"reflect"
	"strconv"
)

type Traverser struct {
	OriJson  map[string]interface{}
	Selector []*Selector
}

func NewTraverser(oriJson map[string]interface{}, Selector []*Selector) *Traverser {
	return &Traverser{OriJson: oriJson, Selector: Selector}
}
func GetValue(json map[string]interface{}, selector []*Selector) (interface{}, error) {
	// use iterative traversal so we can save stack memory here
	curJson := json
	for idx := 0; idx < len(selector); idx++ {
		val := selector[idx]
		if val.Token == IDENTIFIER {
			if _, ok := curJson[val.FieldName]; ok {
				dd := curJson[val.FieldName]
				rt := reflect.TypeOf(dd)
				switch rt.Kind() {
				case reflect.Map:
					if idx == len(selector) {
						return dd, nil
					} else {
						curJson = dd.(map[string]interface{})
						continue
					}
				case reflect.Slice:
					if idx == len(selector) {
						return dd, nil
					}
					// TODO: handling arrays

					if idx+4 == len(selector) {
						// return specific index
						sliceIndex, err := strconv.Atoi(selector[idx+2].Number)
						if err != nil {
							return nil, err
						}
						rv := reflect.ValueOf(dd)

						SubData := rv.Index(sliceIndex)
						return SubData.Interface(), nil
					} else {
						// go deeper
						sliceIndex, err := strconv.Atoi(selector[idx+2].Number)
						if err != nil {
							return nil, err
						}
						rv := reflect.ValueOf(dd)

						SubData := rv.Index(sliceIndex)
						idx += 3
						curJson = SubData.Interface().(map[string]interface{})
						continue
					}
					// return nil, nil
				default:
					return dd, nil
				}
			}
		} else {
			// this is evaluable expession.
			// TODO: implement aggregate function here
			if isAggregate(val.Token) {
				aggregateFunction := val.Token
				allParam := []*Selector{}
				dollarSign := &Selector{Token: DOLLARSIGN}
				allParam = append(allParam, dollarSign)
				allParam = append(allParam, val.FunctionParam...)
				subTraverser := NewTraverser(curJson, allParam)
				rawArray, err := subTraverser.Traverse()
				if err != nil {
					return nil, err
				}
				switch aggregateFunction {
				case SUM:
					// fmt.Println("SUM on", rawArray)
					return Sum(rawArray), nil
				case AVG:
					//fmt.Println("AVG on", rawArray)
					return Avg(rawArray), nil
				}
			}
			return nil, nil
		}
	}
	return curJson, nil
}
func (t *Traverser) Traverse() (interface{}, error) {
	if len(t.Selector) > 0 {
		// traverse this json
		if t.Selector[0].Token == DOLLARSIGN {
			if t.Selector[1].Token == IDENTIFIER {
				fieldName := t.Selector[1].FieldName
				data := t.OriJson[fieldName]
				if _, ok := data.(map[string]interface{}); ok {
					// we can go deeper
					if len(t.Selector) == 2 {
						// not going deeper because we only interested in level1
						return data, nil
					} else {
						// we go deeper
						return GetValue(data.(map[string]interface{}), t.Selector[2:])
					}

				} else if rt := reflect.TypeOf(data); rt.Kind() == reflect.Slice {
					// data is array
					// fmt.Println("Data is array")
					if len(t.Selector) >= 5 {
						// check if we have index selector
						index, err := strconv.Atoi(t.Selector[3].Number)
						if err != nil {
							return nil, err
						}
						rv := reflect.ValueOf(data)

						SubData := rv.Index(index)

						if len(t.Selector) == 5 {
							return SubData.Interface(), nil
						} else {
							// fmt.Println(SubData)
							// go deeper
							return GetValue(SubData.Interface().(map[string]interface{}), t.Selector[5:])
						}
					} else {
						// return as is
						return data, nil
					}
				} else {
					// data is primitive
					return data, nil
				}
			} else {
				if isAggregate(t.Selector[1].Token) {
					// handle aggregate here
					aggregateFunction := t.Selector[1].Token
					allParam := []*Selector{}
					dollarSign := &Selector{Token: DOLLARSIGN}
					allParam = append(allParam, dollarSign)
					allParam = append(allParam, t.Selector[1].FunctionParam...)
					subTraverser := NewTraverser(t.OriJson, allParam)
					rawArray, err := subTraverser.Traverse()
					if err != nil {
						return nil, err
					}
					switch aggregateFunction {
					case SUM:
						// fmt.Println("SUM on", rawArray)
						return Sum(rawArray), nil
					case AVG:
						// fmt.Println("AVG on", rawArray)
						return Avg(rawArray), nil
					}
				} else if t.Selector[1].Token == SQUAREBRACKET_OPEN {
					// its an array
					if t.Selector[2].Token == NUMBER {
						_, err := strconv.Atoi(t.Selector[2].Number)
						if err != nil {
							return nil, err
						}

					}
				}
			}
		} else {
			// this means a static field. Evaluate this
		}
	}
	return nil, nil
}
func Avg(param interface{}) float64 {
	result := 0.0
	rv := reflect.ValueOf(param)
	count := rv.Len()
	if rv.Kind() == reflect.Slice {
		for i := 0; i < count; i++ {
			value := rv.Index(i)
			if value.CanFloat() {
				result += value.Float()
			} else if value.CanInt() {
				result += float64(value.Int())
			}
		}
	}
	return result / float64(count)
}
func Sum(param interface{}) float64 {
	result := 0.0
	rv := reflect.ValueOf(param)
	if rv.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			value := rv.Index(i)
			if value.CanFloat() {
				result += value.Float()
			} else if value.CanInt() {
				result += float64(value.Int())
			}
		}
	}
	return result
}
