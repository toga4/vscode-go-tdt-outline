package main_test

import "testing"

func TestNestedTable(t *testing.T) {
	type Want struct {
		Name string
	}
	tests := []struct {
		name  string
		input int
		want  []Want
	}{
		{
			name:  "normal case",
			input: 1,
			want: []Want{
				{
					Name: "name1",
				},
			},
		},
		{
			name:  "zero value",
			input: 0,
			want: []Want{
				{
					Name: "name2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// test logic
		})
	}
}
