package main_test

import "testing"

func TestNoName(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{input: 1, want: 1},
		{input: 2, want: 2},
	}
}
