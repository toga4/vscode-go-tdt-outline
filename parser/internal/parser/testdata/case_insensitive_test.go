package main_test

import "testing"

func TestCaseInsensitive(t *testing.T) {
	tests := []struct {
		NAME string
		Name string
		name string
	}{
		{NAME: "大文字NAME"},
		{Name: "混合Name"},
		{name: "小文字name"},
	}
}
