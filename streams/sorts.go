package streams

import (
	"strings"
)

func SortStringsAsc() SortFunc {
	return func(i interface{}, j interface{}) int {
		return strings.Compare(i.(string), j.(string))
	}
}

func SortStringsDesc() SortFunc {
	return func(i interface{}, j interface{}) int {
		return strings.Compare(j.(string), i.(string))
	}
}
