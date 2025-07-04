package main_test

import "testing"

func TestDeclaration(t *testing.T) {
	var tests = []struct {
		name  string
		input int
		want  int
	}{
		{
			name:  "normal case",
			input: 1,
			want:  1,
		},
		{
			name:  "zero value",
			input: 0,
			want:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// test logic
		})
	}
}
