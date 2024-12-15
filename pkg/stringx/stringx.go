package stringx

import (
	"regexp"
	"slices"
	"strings"
)

// FilterAlphanumeric filter out special characters and retain only letters and numbers in the string
func FilterAlphanumeric(input string) string {
	re := regexp.MustCompile("[^a-zA-Z0-9]")
	return re.ReplaceAllString(input, "")
}

func FuzzyContains(slice []string, target string) bool {
	return slices.ContainsFunc(slice, func(s string) bool {
		return strings.Contains(target, s)
	})
}
