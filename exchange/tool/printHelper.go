package tool

import (
	"fmt"
	"reflect"
)

func PrintStructDetail(p interface{}) {
	val := reflect.ValueOf(p)
	typ := reflect.TypeOf(p)

	for i := 0; i < val.NumField(); i++ {
		fmt.Printf("%2d,%10s, %s:%v", i+1, val.Field(i).Type(), typ.Field(i).Name, val.Field(i).Interface())
		fmt.Println("")
	}
}

func PrintStructNameValue(p interface{}) {
	val := reflect.ValueOf(p)
	typ := reflect.TypeOf(p)

	for i := 0; i < val.NumField(); i++ {
		fmt.Printf("%s:%v,", typ.Field(i).Name, val.Field(i).Interface())
	}
	fmt.Println("")
}
