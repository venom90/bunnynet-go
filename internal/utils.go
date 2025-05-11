package internal

import (
	"strconv"
)

// FormatInt64 converts an int64 to a string
func FormatInt64(i int64) string {
	return strconv.FormatInt(i, 10)
}
