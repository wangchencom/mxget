package utils

import (
	"regexp"
	"strings"
)

func TrimInvalidFilePathChars(path string) string {
	re := regexp.MustCompile(`[\\/:*?"<>|]`)
	path = re.ReplaceAllString(path, " ")
	return strings.TrimSpace(path)
}
