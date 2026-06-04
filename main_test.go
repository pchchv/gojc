package main

import (
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	go server()

	// Give the server 500 milliseconds to guarantee that it will come up and take over the port
	time.Sleep(500 * time.Millisecond)

	// Run all tests of the package
	exitCode := m.Run()

	// Wrap up the process with the test exit code
	os.Exit(exitCode)
}
