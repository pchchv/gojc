package main

import (
	"log"
	"os"

	"github.com/pchchv/env"
)

var filename string

func init() {
	// Load values from .env into the system
	if err := env.Load(); err != nil {
		log.Panic("no .env file found")
	}

	filename = getEnvValue("FILENAME")
}

func getEnvValue(v string) string {
	// Getting a value
	// Outputs a panic if the value is missing
	value, exist := os.LookupEnv(v)
	if !exist {
		log.Panicf("value %v does not exist", v)
	}
	return value
}

func main() {}
