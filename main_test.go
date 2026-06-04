package main

import (
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Set default environment variables if they are not specified in CI
	if os.Getenv("HOST") == "" {
		os.Setenv("HOST", "127.0.0.1")
	}
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8080")
	}
	if os.Getenv("FILENAME") == "" {
		os.Setenv("FILENAME", "config.json")
	}

	go server()

	// Give the server 500 milliseconds to guarantee that it will come up and take over the port
	time.Sleep(500 * time.Millisecond)

	// Run all tests of the package
	exitCode := m.Run()

	// Wrap up the process with the test exit code
	os.Exit(exitCode)
}
