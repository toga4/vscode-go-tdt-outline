package main

import "testing"

func Example() {
	tests := []struct {
		name string
	}{
		{name: "should be ignored"},
	}
}

func BenchmarkExample(b *testing.B) {
	tests := []struct {
		name string
	}{
		{name: "should be ignored"},
	}
}

func TestNonTestFunction(t *testing.T)
