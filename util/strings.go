package util

import (
	"regexp"
	"strings"
)

var re = regexp.MustCompile("[^a-z0-9]+")

func ConvertToSlug(input string) string {
	return strings.Trim(re.ReplaceAllString(strings.ToLower(input), "-"), "-")
}
