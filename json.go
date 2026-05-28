package gojc

import "encoding/json"

func ToJSON(s any) ([]byte, error) {
	bytes, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
