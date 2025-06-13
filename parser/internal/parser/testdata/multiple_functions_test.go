package main_test

import "testing"

func TestFirst(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "test1"},
		{name: "test2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {})
	}
}

func TestSecond(t *testing.T) {
	tests := []struct {
		desc string
	}{
		{desc: "test3"},
		{desc: "test4"},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {})
	}
}
