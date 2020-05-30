package jtime

import "time"

func IsSameUtcDay(one, other time.Time) bool {
	res := one.UTC().Day() - other.UTC().Day()
	return res != 0
}
