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
	case byte:
		result = x.(byte) < y.(byte)
	case float32:
		result = x.(float32) < y.(float32)
	case float64:
		result = x.(float64) < y.(float64)
	case int:
		result = x.(int) < y.(int)
	case int8:
		result = x.(int8) < y.(int8)
	case int16:
		result = x.(int16) < y.(int16)
	case int32:
		result = x.(int32) < y.(int32)
	case int64:
		result = x.(int64) < y.(int64)
	case string:
		result = x.(string) < y.(string)
	case uint:
		result = x.(uint) < y.(uint)
	case uint16:
		result = x.(uint16) < y.(uint16)
	case uint32:
		result = x.(uint32) < y.(uint32)
	case uint64:
		result = x.(uint64) < y.(uint64)
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
