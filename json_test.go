package gojc

import (
	"bytes"
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
				Location: "New York", // Field tagged with "-" must be omitted
			},
			want:    []byte(`{"title":"Acme Corp"}`),
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
