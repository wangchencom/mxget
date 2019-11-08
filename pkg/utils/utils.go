package utils

import (
	"regexp"
	"strings"
)

var (
	re = regexp.MustCompile(`[\\/:*?"<>|]`)
)

func TrimInvalidFilePathChars(path string) string {
	path = re.ReplaceAllString(path, " ")
	return strings.TrimSpace(path)
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
