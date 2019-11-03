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
