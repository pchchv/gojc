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
			name:    "Successful serialization of a simple structure",
			input:   User{Id: 111, Name: "John", Email: "john@company.com", IsActive: true},
			want:    []byte(`{"id":111,"name":"John","email":"john@company.com","isActive":true}`),
			wantErr: false,
		},
		{
			name: "Ignoring hidden fields and omitempty",
			input: Company{
				Title:    "TechCorp",
				CEO:      nil,
				Location: "Moscow",
			},
			want:    []byte(`{"title":"TechCorp"}`),
			wantErr: false,
		},
		{
			name:    "Passing an empty interface (nil)",
			input:   nil,
			want:    []byte(`null`),
			wantErr: false,
		},
		{
			name:    "Error thrown (invalid type for JSON - channel)",
			input:   make(chan int),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ToJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("ToJSON() = %s, want %s", string(got), string(tt.want))
			}
		})
	}
}
