package main_test

import "testing"

func TestPositionalFieldForm(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  int
	}{
		{"normal case", 1, 1},
		{"zero value", 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// test logic
		})
	}
}
