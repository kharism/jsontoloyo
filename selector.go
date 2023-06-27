package jsontoloyo

type Selector struct {
	Token         Token
	Number        string
	FieldName     string
	FunctionName  string
	FunctionParam []*Selector
}

func ArrSelectorToStr(params []*Selector) string {
	output := ""
	for _, v := range params {
		output += v.String()
	}
	return output
}
func (s *Selector) String() string {
	if s.FieldName != "" {
		return s.FieldName
	} else if s.FunctionName != "" {
		if s.FunctionParam != nil {
			return "$" + s.FunctionName + "(" + ArrSelectorToStr(s.FunctionParam) + ")"
		} else {
			return "$" + s.FunctionName + "()"
		}

	} else if s.Number != "" {
		return s.Number
	}
	return ""
}
