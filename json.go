package gojc

import "encoding/json"

// ToJSON converts any struct into a native JSON byte slice.
func ToJSON(v any) ([]byte, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
