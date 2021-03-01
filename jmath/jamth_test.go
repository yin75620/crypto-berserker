package jmath

import (
	"fmt"
	"testing"
)

func TestMath(t *testing.T) {

	a := FloatFloorByFloatCount(3.127756, 10.003)
	fmt.Println(a)
	if a != 3.127 {
		t.Error()
	}
}
