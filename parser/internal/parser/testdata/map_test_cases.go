package main_test

import "testing"

func TestWithMap(t *testing.T) {
	tests := map[string]struct {
		input    int
		expected int
		wantErr  bool
	}{
		"正常系: 基本的なケース": {
			input:    1,
			expected: 2,
			wantErr:  false,
		},
		"正常系: ゼロ値": {
			input:    0,
			expected: 0,
			wantErr:  false,
		},
		"異常系: 負の値": {
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

// シンプルなmap定義
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

// 型定義を使用したmap
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
			input: "こんにちは",
			want:  "こんにちは",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// test logic
			_ = tc
		})
	}
}
