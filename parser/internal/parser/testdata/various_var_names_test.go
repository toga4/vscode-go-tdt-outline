package main_test

import "testing"

func TestScenarios(t *testing.T) {
	scenarios := []struct {
		name string
		data string
	}{
		{name: "シナリオ1"},
		{name: "シナリオ2"},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// test logic
		})
	}
}

func TestExamples(t *testing.T) {
	examples := []struct {
		desc string
		code string
	}{
		{desc: "例1: 基本"},
		{desc: "例2: 応用"},
	}
	for _, ex := range examples {
		t.Run(ex.desc, func(t *testing.T) {
			// test logic
		})
	}
}
