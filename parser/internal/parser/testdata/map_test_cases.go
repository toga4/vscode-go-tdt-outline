package main_test

import "testing"

func TestWithMap(t *testing.T) {
	tests := map[string]struct {
		input    int
		expected int
		wantErr  bool
	}{
		"normal case: basic scenario": {
			input:    1,
			expected: 2,
			wantErr:  false,
		},
		"normal case: zero value": {
			input:    0,
			expected: 0,
			wantErr:  false,
		},
		"error case: negative value": {
			input:    -1,
			expected: -1,
			wantErr:  true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// test logic
		})
	}
}

// Simple map definition
func TestSimpleMap(t *testing.T) {
	cases := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for name, value := range cases {
		t.Run(name, func(t *testing.T) {
			// test logic with value
			_ = value
		})
	}
}

// Map using type definition
type TestCase struct {
	input string
	want  string
}

func TestTypedMap(t *testing.T) {
	tests := map[string]TestCase{
		"empty string": {
			input: "",
			want:  "",
		},
		"hello world": {
			input: "hello",
			want:  "HELLO",
		},
		"unicode": {
			input: "hello world",
			want:  "HELLO WORLD",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// test logic
			_ = tc
		})
	}
}
