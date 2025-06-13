package main_test

import "testing"

type Test struct {
	name    string
	input   int
	want    int
	wantErr bool
}

func TestTypedStruct(t *testing.T) {
	tests := []Test{
		{
			name:    "normal case: basic scenario",
			input:   1,
			want:    2,
			wantErr: false,
		},
		{
			name:    "normal case: zero value scenario",
			input:   0,
			want:    0,
			wantErr: false,
		},
		{
			name:    "error case: invalid input",
			input:   -1,
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ... test logic ...
		})
	}
}

type Tests []Test

func TestTypeAlias(t *testing.T) {
	tests := Tests{
		{
			name:    "type alias: case 1",
			input:   1,
			want:    2,
			wantErr: false,
		},
		{
			name:    "type alias: case 2",
			input:   0,
			want:    0,
			wantErr: false,
		},
		{
			name:    "type alias: case 3",
			input:   -1,
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ... test logic ...
		})
	}
}
