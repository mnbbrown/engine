package engine

import (
	"fmt"
	"math"
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func humanSize(size int) string {
	suffix := "B"
	postfixes := []string{"", "K", "M", "G", "T", "P", "E", "Z"}
	for _, v := range postfixes {
		if math.Abs(float64(size)) < 1024.0 {
			return fmt.Sprintf("%d %s%s", size, v, suffix)
		}
		size /= 1024.0
	}
	return fmt.Sprintf("%d %s%s", size, "Y", suffix)
}