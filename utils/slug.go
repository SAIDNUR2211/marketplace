package utils

import (
	"regexp"
	"strings"
)

var (
	nonAlphaNumeric  = regexp.MustCompile(`[^a-z0-9-]+`)
	duplicateHyphens = regexp.MustCompile(`-+`)
)

func GenerateSlug(text string) string {
	slug := strings.ToLower(text)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = nonAlphaNumeric.ReplaceAllString(slug, "")
	slug = duplicateHyphens.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}
