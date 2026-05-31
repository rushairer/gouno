package utility

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var camelCaseRegex = regexp.MustCompile("([a-z0-9])([A-Z])")

// ToCamelCase converts a snake_case string to PascalCase.
// For example: "hello_world" becomes "HelloWorld".
func ToCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		parts[i] = cases.Title(language.English).String(part)
	}
	return strings.Join(parts, "")
}

// ToSnakeCase converts a camelCase or PascalCase string to snake_case.
// For example: "HelloWorld" becomes "hello_world".
func ToSnakeCase(s string) string {
	s = camelCaseRegex.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(s)
}
