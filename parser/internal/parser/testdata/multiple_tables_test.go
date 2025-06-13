package main_test

import "testing"

func TestMultipleTables(t *testing.T) {
	// First test table
	tests1 := []struct {
		name string
	}{
		{name: "table1-test1"},
		{name: "table1-test2"},
	}

	// Second test table
	tests2 := []struct {
		name string
	}{
		{name: "table2-test1"},
		{name: "table2-test2"},
	}
}
