package main

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

// SaveJSONToFile appends a new JSON object into a valid, strict JSON array file.
func SaveJSONToFile(data []byte) error {
	fileInfo, err := os.Stat(filename)

	// If the file does not exist or is empty, initialize it as a fresh array
	if os.IsNotExist(err) || fileInfo.Size() == 0 {
		// Wrap payload inside an array structure: [data]
		payload := append(append([]byte("["), data...), []byte("]")...)
		return os.WriteFile(filename, payload, 0644)
	} else if err != nil {
		return err
	}

	// Open the existing file for read/write modifications
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Seek to the position right before the closing bracket ']'
	// We subtract 1 from the end of the file
	_, err = file.Seek(-1, 2)
	if err != nil {
		return err
	}

	// Construct the append payload: a comma delimiter followed by our new data and closing array bracket
	payload := append(append([]byte(","), data...), []byte("]")...)

	_, err = file.Write(payload)
	return err
}
