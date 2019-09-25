package jmath

import (
	"math"
)

func FloatFloorByFloat(f float64, unit float64) float64 {
	if unit == 0 {
		return f //不運算
	}
	val := f / unit
	val = math.Floor(val)
	res := val * unit
	return res
}

func FloatFloor(f float64, count int) float64 {
	val := f * math.Pow10(count)
	val = math.Floor(val)
	res := val / math.Pow10(count)
	return res
}
