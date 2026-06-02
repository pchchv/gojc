package main

import "strings"

// exportedName capitalizes the first letter (required for struct fields in Go).
func exportedName(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
