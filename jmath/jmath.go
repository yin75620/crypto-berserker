package jmath

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// unit is like 3.123456 or 125.1234
func FloatFloorByFloatCount(f float64, unit float64) float64 {
	if unit == 0 {
		return f //不運算
	}
	str := strconv.FormatFloat(unit, 'f', -1, 64)
	index := strings.LastIndex(str, ".")
	floatcount := len(str) - 1 - index

	return FloatFloor(f, floatcount)
}

// unit is mean 0.001 or 0.0001 and so
func FloatFloorByFloat(f float64, unit float64) float64 {
	if unit == 0 {
		return f //不運算
	}
	val := f / unit
	val = math.Floor(val)
	res := val * unit

	count := int(-math.Floor(math.Log10(unit)))

	floatStr := fmt.Sprintf("%."+strconv.Itoa(count)+"f", res)
	inst, _ := strconv.ParseFloat(floatStr, 64)
	return inst
}

func FloatFloor(f float64, count int) float64 {
	val := f * math.Pow10(count)
	val = math.Floor(val)
	res := val / math.Pow10(count)
	return res
}
