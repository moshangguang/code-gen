package strutils

import "strings"

func PrefixPattern(data []string, pattern string) []string {
	pattern = strings.TrimSpace(pattern)
	if len(pattern) == 0 {
		return data
	}
	result := make([]string, 0, len(data)/2)
	for _, d := range data {
		if strings.HasPrefix(d, pattern) {
			result = append(result, d)
		}
	}
	return result
}

func PatternTrimSpace(data []string, pattern string) []string {
	pattern = strings.TrimSpace(pattern)
	if len(pattern) == 0 {
		return data
	}
	result := make([]string, 0, len(data)/2)
	for _, d := range data {
		if strings.Contains(d, pattern) {
			result = append(result, d)
		}
	}
	return result
}
func IsEmptyString(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}
func IsNotEmptyString(str string) bool {
	return !IsEmptyString(str)
}
