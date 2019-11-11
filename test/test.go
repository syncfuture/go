package test

import (
	"math/rand"
	"reflect"

	r "github.com/syncfuture/go/rand"
)

func AutoFixture(objPtr interface{}) {
	t := reflect.TypeOf(objPtr)
	if t.Kind() != reflect.Ptr {
		panic("objPtr must be pointer")
	}

	v := reflect.ValueOf(objPtr).Elem()

	fieldCount := v.NumField()
	for i := 0; i < fieldCount; i++ {
		field := v.Field(i)
		if field.CanAddr() {
			switch field.Kind() {
			case reflect.String:
				field.SetString(r.String(10))
			case reflect.Int:
				field.SetInt(rand.Int63())
			case reflect.Int32:
				field.SetInt(rand.Int63())
			case reflect.Int64:
				field.SetInt(rand.Int63())
			case reflect.Int16:
				field.SetInt(rand.Int63())
			case reflect.Int8:
				field.SetInt(rand.Int63())
			case reflect.Float32:
				field.SetFloat(rand.Float64())
			case reflect.Float64:
				field.SetFloat(rand.Float64())
			}
		}
	}
}
