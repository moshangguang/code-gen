package help

import "strings"

func IsEmptyString(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}
func IsNotEmptyString(str string) bool {
	return !IsEmptyString(str)
}
