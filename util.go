package goroutine

import (
	"encoding/json"
	"fmt"
	"runtime"
)

func SpotHere() string {
	_, f, l, _ := runtime.Caller(2)
	return fmt.Sprintf("%s:%d", f, l)
}

func SpotHereV2(depth int) string {
	_, f, l, _ := runtime.Caller(depth)
	return fmt.Sprintf("%s:%d", f, l)
}


func JSON(src interface{}) string {
	b, _ := json.MarshalIndent(src, "", "  ")
	return string(b)
}

func DefaultInt(src int, defaultv int) int {
	if src == 0 {
		return defaultv
	}
	return src
}
