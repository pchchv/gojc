package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	IsActive bool   `json:"isActive"`
}

type Company struct {
	Title    string `json:"title"`
	CEO      *User  `json:"ceo,omitempty"`
	Location string `json:"-"`
}

func TestToJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    []byte
		wantErr bool
	}{
		{
			name: "Successfully marshal basic User struct",
			input: User{
				Id:       1,
				Name:     "John Doe",
				Email:    "john@example.com",
				IsActive: true,
			},
			want:    []byte(`{"id":1,"name":"John Doe","email":"john@example.com","isActive":true}`),
			wantErr: false,
		},
		{
			name: "Handle Company struct with empty CEO (omitempty) and ignored field (-)",
			input: Company{
				Title:    "Acme Corp",
				CEO:      nil,
				Location: "New York",
			},
			want:    []byte(`{"title":"Acme Corp"}`),
			wantErr: false,
		},
		{
			name: "Handle Company struct with nested CEO populated",
			input: Company{
				Title: "Tech Solutions",
				CEO: &User{
					Id:       2,
					Name:     "Alice Smith",
					Email:    "alice@example.com",
					IsActive: true,
				},
				Location: "London",
			},
			want:    []byte(`{"title":"Tech Solutions","ceo":{"id":2,"name":"Alice Smith","email":"alice@example.com","isActive":true}}`),
			wantErr: false,
		},
		{
			name:    "Passing an empty interface (nil)",
			input:   nil,
			want:    []byte(`null`),
			wantErr: false,
		},
		{
			name:    "Return error for types that cannot be serialized (channel)",
			input:   make(chan int),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ToJSON() unexpected error status: got error = %v, wantErr = %v", err, tt.wantErr)
			}

			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("ToJSON() mismatch.\ngot : %s\nwant: %s", string(got), string(tt.want))
			}
		})
	}
}

func TestSaveJSONToFile(t *testing.T) {
	userPayload := []byte(`{"id":1,"name":"John Doe","email":"john@example.com","isActive":true}`)
	companyPayload := []byte(`{"title":"Acme Corp"}`)

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "array_output.json")

	// Mock the environment variable safely for this test execution scope
	t.Setenv("FILENAME", filePath)

	// Preserve the original package global state and restore it when this test finishes
	oldFilename := filename
	defer func() { filename = oldFilename }()

	// Mimic the production initialization logic so the global variable holds the target path
	// In your real main app, this would be: filename = getEnvValue("FILENAME")
	filename = filePath
	t.Run("Append objects into a strictly valid JSON array using global configuration", func(t *testing.T) {
		// First append: Creates the file and initializes the array
		err := SaveJSONToFile(userPayload)
		if err != nil {
			t.Fatalf("SaveJSONToFile() initial write failed: %v", err)
		}

		// Second append: Inserts comma and updates the closing bracket
		err = SaveJSONToFile(companyPayload)
		if err != nil {
			t.Fatalf("SaveJSONToFile() second append failed: %v", err)
		}

		// Read back results
		readData, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read array data file: %v", err)
		}

		// Expected format is a valid, un-nested standard JSON array
		expected := []byte(`[{"id":1,"name":"John Doe","email":"john@example.com","isActive":true},{"title":"Acme Corp"}]`)

		if !bytes.Equal(readData, expected) {
			t.Errorf("Strict JSON array mismatch.\nGot : %s\nWant: %s", string(readData), string(expected))
		}
	})
}
