package ginx

import (
	"fmt"
	"strings"
)

func MergeURL(first, second string) string {
	first = strings.Trim(first, "/")
	second = strings.Trim(second, "/")
	if first == "" && second == "" {
		return "/"
	} else if first == "" {
		return fmt.Sprintf("/%s", second)
	} else if second == "" {
		return fmt.Sprintf("/%s", first)
	}
	return fmt.Sprintf("/%s/%s", first, second)
}
