package helper

import (
	"strings"
)

func Equals[T ~string](a, b T) bool {
	return strings.EqualFold(string(a), string(b))
}
