package debug_tool

import (
	"fmt"
)

const Dev = true

func Debug(a ...interface{}) {
	fmt.Println(a...)
}

// 可參考此網站進行優化
// https://mzh.io/golang-build-tags-for-debug/
