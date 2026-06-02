package main

import (
	"encoding/json"
	"strconv"
	"strings"
)

// parseAnyType attempts to cast a string to the closest Go type possible.
func parseAnyType(val string) any {
	// Check if it's nested JSON (array or object)
	// For example: [1,2,3] or {"city":"Moscow"}
	if (strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]")) ||
		(strings.HasPrefix(val, "{") && strings.HasSuffix(val, "}")) {
		var jsonRaw any
		if err := json.Unmarshal([]byte(val), &jsonRaw); err == nil {
			// Return []any or map[string]any
			return jsonRaw
		}
	}

	// Try to convert it to an integer (int)
	if i, err := strconv.Atoi(val); err == nil {
		return i
	}

	// Try converting to a floating point number (float64)
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		return f
	}

	// Try converting to a bool
	if b, err := strconv.ParseBool(val); err == nil {
		return b
	}

	// If nothing matches, return it as a string
	return val
}

// exportedName capitalizes the first letter (required for struct fields in Go).
func exportedName(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
