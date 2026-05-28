package gojc

import (
	"encoding/json"
	"os"
)

// ToJSON converts any struct into a native JSON byte slice.
func ToJSON(v any) ([]byte, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// SaveJSONToFile writes the JSON byte slice into a file with the specified name.
func SaveJSONToFile(filename string, data []byte) error {
	// Open file with flags:
	//   Append
	//   Create if non-existent
	//   Write-Only
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the JSON data payload
	if _, err := file.Write(data); err != nil {
		return err
	}

	// Append a newline character to delimit multiple JSON objects properly
	if _, err := file.Write([]byte("\n")); err != nil {
		return err
	}

	return nil
}
