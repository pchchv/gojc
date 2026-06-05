package main

import (
	"log"
	"os"
	"testing"
	"time"
)

// setupDummyEnv runs before any init() function in the package.
// It ensures a physical .env file exists so env.Load() doesn't panic.
var _ = func() bool {
	err := os.WriteFile(".env", []byte("HOST=127.0.0.1\nPORT=8080\nFILENAME=test_collect.json\n"), 0644)
	if err != nil {
		log.Fatalf("Failed to create temporary .env for testing: %v", err)
	}
	return true
}()

func TestMain(m *testing.M) {
	// Clean up the temporary .env file after all tests finish executing
	defer func() {
		_ = os.Remove(".env")
	}()

	log.Println("[TestMain] Starting test API server...")

	go server()

	// Give the server 500 milliseconds to guarantee that it will come up and take over the port
	time.Sleep(500 * time.Millisecond)

	// Run all tests of the package
	exitCode := m.Run()

	// Wrap up the process with the test exit code
	os.Exit(exitCode)
}
