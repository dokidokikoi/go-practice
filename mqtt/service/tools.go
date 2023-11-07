package service

import "reflect"

func IsZero(v interface{}) bool {
	rv := reflect.ValueOf(v)
	zv := reflect.Zero(rv.Type())

	return reflect.DeepEqual(rv.Interface(), zv.Interface())
}
