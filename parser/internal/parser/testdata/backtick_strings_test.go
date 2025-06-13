package main

import "testing"

func TestBacktickStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "double quote string",
			input:    "test",
			expected: "result",
		},
		{
			name:     `backtick string`,
			input:    "test",
			expected: "result",
		},
		{
			name:     `backtick with "quotes"`,
			input:    "test",
			expected: "result",
		},
		{
			name: `backtick with
newlines`,
			input:    "test",
			expected: "result",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// test logic here
		})
	}
}
