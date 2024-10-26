package jtime

import "time"

func IsSameUtcDay(one, other time.Time) bool {
	res := one.UTC().Day() - other.UTC().Day()
	return res != 0
}

func UnixToTime(unix int64) time.Time {
	unimilli := unix

	// 將毫秒轉為秒和奈秒
	seconds := unimilli / 1000
	nanoseconds := (unimilli % 1000) * 1e6

	// 使用 time.Unix 函數轉換為 time.Time
	return time.Unix(seconds, nanoseconds)
}
