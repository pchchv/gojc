package main

import (
	"fmt"
	"log"

	"github.com/pchchv/gojc"
)

type Config struct {
	AppID   string `json:"app_id"`
	Version string `json:"version"`
}

func main() {
	cfg := Config{AppID: "my_app", Version: "1.0.0"}

	// Converting to a byte slice
	data, err := gojc.ToJSON(cfg)
	if err != nil {
		log.Fatalf("serialization error: %v", err)
	}

	// Save to file
	err = gojc.SaveJSONToFile(data)
	if err != nil {
		log.Fatalf("error writing to file: %v", err)
	}

	fmt.Println("The config.json file has been saved successfully.!")
}
