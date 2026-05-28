package gojc

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

func TestAppendJSONToFile(t *testing.T) {
	// Prepare test data payloads
	userPayload := []byte(`{"id":1,"name":"John Doe","email":"john@example.com","isActive":true}`)
	companyPayload := []byte(`{"title":"Acme Corp"}`)

	// Create an isolated temporary directory for testing filesystem interactions
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "stream_output.json")

	t.Run("Append multiple JSON objects successfully", func(t *testing.T) {
		// First entry: writes to a brand new file
		err := SaveJSONToFile(filePath, userPayload)
		if err != nil {
			t.Fatalf("AppendJSONToFile() first write failed: %v", err)
		}

		// Second entry: appends to the existing file
		err = SaveJSONToFile(filePath, companyPayload)
		if err != nil {
			t.Fatalf("AppendJSONToFile() second append failed: %v", err)
		}

		// Read back the entire accumulated file content
		readData, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read the appended data file: %v", err)
		}

		// Construct the expected outcome (both entries, each followed by a newline)
		var expectedBuffer bytes.Buffer
		expectedBuffer.Write(userPayload)
		expectedBuffer.WriteByte('\n')
		expectedBuffer.Write(companyPayload)
		expectedBuffer.WriteByte('\n')
		expected := expectedBuffer.Bytes()

		if !bytes.Equal(readData, expected) {
			t.Errorf("Appended data mismatch.\nGot contents:\n%s\nExpected contents:\n%s", string(readData), string(expected))
		}
	})

	t.Run("Fail when appending to a missing path hierarchy", func(t *testing.T) {
		invalidPath := filepath.Join(tempDir, "missing_folder", "output.json")
		err := SaveJSONToFile(invalidPath, userPayload)
		if err == nil {
			t.Error("Expected a filesystem error due to non-existent directory chain, but got nil")
		}
	})
}
