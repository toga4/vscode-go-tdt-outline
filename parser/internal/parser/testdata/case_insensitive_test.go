package main_test

import "testing"

func TestCaseInsensitive(t *testing.T) {
	tests := []struct {
		NAME string
		Name string
		name string
	}{
		{NAME: "uppercase NAME"},
		{Name: "mixed case Name"},
		{name: "lowercase name"},
	}
}
