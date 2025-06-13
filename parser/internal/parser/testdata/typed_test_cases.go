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
			name:    "正常系: 基本的なケース",
			input:   1,
			want:    2,
			wantErr: false,
		},
		{
			name:    "正常系: ゼロ値のケース",
			input:   0,
			want:    0,
			wantErr: false,
		},
		{
			name:    "異常系: 不正な入力",
			input:   -1,
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ... テストロジック ...
		})
	}
}

type Tests []Test

func TestTypeAlias(t *testing.T) {
	tests := Tests{
		{
			name:    "型エイリアス: ケース1",
			input:   1,
			want:    2,
			wantErr: false,
		},
		{
			name:    "型エイリアス: ケース2",
			input:   0,
			want:    0,
			wantErr: false,
		},
		{
			name:    "型エイリアス: ケース3",
			input:   -1,
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ... テストロジック ...
		})
	}
}
