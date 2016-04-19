package lmaodb

import "reflect"

type sortable struct {
	field      string
	sliceValue reflect.Value
}

func newSortable(slice interface{}, field string) *sortable {
	s := sortable{
		field:      field,
		sliceValue: reflect.ValueOf(slice),
	}
	return &s
}

func (s *sortable) Len() int {
	return s.sliceValue.Len()
}

func (s *sortable) Less(i, j int) (result bool) {
	x := s.sliceValue.Index(i).FieldByName(s.field).Interface()
	y := s.sliceValue.Index(j).FieldByName(s.field).Interface()
	switch x.(type) {
	case int:
		result = x.(int) < y.(int)
	case int64:
		result = x.(int64) < y.(int64)
	case string:
		result = x.(string) < y.(string)
	}
	return
}

func (s *sortable) Swap(i, j int) {
	x := s.sliceValue.Index(i)
	y := s.sliceValue.Index(j)
	temp := reflect.Indirect(reflect.New(y.Type()))
	temp.Set(x)
	x.Set(y)
	y.Set(temp)
}
