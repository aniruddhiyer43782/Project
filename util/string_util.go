package util

import (
	"strings"
)

// EscapeShell escapes shell output by replacing double quotes with single quotes.
func EscapeShell(org string) (dst string) {
	return strings.ReplaceAll(org, "\"", "'")
}
