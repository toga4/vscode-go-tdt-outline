package main_test

import "testing"

func TestVariousFields(t *testing.T) {
	tests := []struct {
		description string
		title       string
		scenario    string
		testName    string
	}{
		{description: "description field"},
		{title: "title field"},
		{scenario: "scenario field"},
		{testName: "testName field"},
	}
}
