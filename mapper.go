package jsontoloyo

import (
	"strings"
)

type Mapper struct {
	FieldName string
	Traverser *Traverser
}

func NewMapperTraverser(fieldname string, traverser *Traverser) *Mapper {
	return &Mapper{FieldName: fieldname, Traverser: traverser}
}
func NewMapperPath(fieldname string, traverserPath string) (*Mapper, error) {
	source := strings.NewReader(traverserPath)
	parser := NewParser(source)
	selectors := []*Selector{}
	err := parser.Parse(&selectors)
	if err != nil {
		return nil, err
	}
	traverser := NewTraverser(nil, selectors)
	return NewMapperTraverser(fieldname, traverser), nil
}

type Mappers struct {
	mapping []*Mapper
}

func NewMappers() *Mappers {
	return &Mappers{mapping: []*Mapper{}}
}
func (m *Mappers) AddMapping(targetFieldname string, path string) (*Mappers, error) {
	newMapper, err := NewMapperPath(targetFieldname, path)
	if err != nil {
		return nil, err
	}
	m.mapping = append(m.mapping, newMapper)
	return m, nil
}
func (m *Mappers) Mapping(source, sink map[string]interface{}) (err error) {
	for _, mapper := range m.mapping {
		mapper.Traverser.OriJson = source
		sink[mapper.FieldName], err = mapper.Traverser.Traverse()
		if err != nil {
			return err
		}
	}
	return nil
}
