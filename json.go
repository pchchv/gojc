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

// SaveJSONToFile writes the JSON byte slice into a
// file with the specified name.
func SaveJSONToFile(filename string, data []byte) error {
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	return nil
}
