package main

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Set default environment variables if they are not specified in CI
	if os.Getenv("HOST") == "" {
		_ = os.Setenv("HOST", "127.0.0.1")
	}
	if os.Getenv("PORT") == "" {
		_ = os.Setenv("PORT", "8080")
	}
	if os.Getenv("FILENAME") == "" {
		_ = os.Setenv("FILENAME", "config.json")
	}

	log.Println("[TestMain] Starting test API server...")

	go server()

	// Give the server 500 milliseconds to guarantee that it will come up and take over the port
	time.Sleep(500 * time.Millisecond)

	// Run all tests of the package
	exitCode := m.Run()

	// Wrap up the process with the test exit code
	os.Exit(exitCode)
}
