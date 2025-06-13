package main_test

import "testing"

func TestExample(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  int
	}{
		{
			name:  "正常系",
			input: 1,
			want:  1,
		},
		{
			name:  "ゼロ値",
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
