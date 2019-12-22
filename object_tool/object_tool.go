package object_tool

import (
	"fmt"
	"log"
)

func ToString(x interface{}) string {
	value := ""
	if w, ok := x.(string); ok {
		value = w
	} else if i, ok := x.(int64); ok {
		value = fmt.Sprintf("%d", i)
	} else if i, ok := x.(int); ok {
		value = fmt.Sprintf("%d", i)
	} else {
		log.Fatal("undefine type:", x)
	}
	return value
}
