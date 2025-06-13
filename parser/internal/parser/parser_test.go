package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     []Symbol
		wantErr  bool
	}{
		{
			name:     "基本的なテーブル駆動テスト",
			filePath: "testdata/basic_table_test.go",
			want: []Symbol{
				{
					Name:   "TestExample",
					Detail: "test function",
					Kind:   SymbolKindFunction,
					Children: []Symbol{
						{
							Name:   "正常系",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "ゼロ値",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "複数のテスト関数",
			filePath: "testdata/multiple_functions_test.go",
			want: []Symbol{
				{
					Name:   "TestFirst",
					Detail: "test function",
					Kind:   SymbolKindFunction,
					Children: []Symbol{
						{
							Name:   "test1",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "test2",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
					},
				},
				{
					Name:   "TestSecond",
					Detail: "test function",
					Kind:   SymbolKindFunction,
					Children: []Symbol{
						{
							Name:   "test3",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "test4",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "異なるフィールド名のサポート",
			filePath: "testdata/various_fields_test.go",
			want: []Symbol{
				{
					Name:   "TestVariousFields",
					Detail: "test function",
					Kind:   SymbolKindFunction,
					Children: []Symbol{
						{
							Name:   "説明フィールド",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "タイトルフィールド",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "シナリオフィールド",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "テスト名フィールド",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "複数のテストテーブル",
			filePath: "testdata/multiple_tables_test.go",
			want: []Symbol{
				{
					Name:   "TestMultipleTables",
					Detail: "test function",
					Kind:   SymbolKindFunction,
					Children: []Symbol{
						{
							Name:   "table1-test1",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "table1-test2",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "table2-test1",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "table2-test2",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "テスト関数ではない関数は無視",
			filePath: "testdata/non_test_functions.go",
			want:     []Symbol{},
			wantErr:  false,
		},
		{
			name:     "nameフィールドがないテストケースは無視",
			filePath: "testdata/no_name_field_test.go",
			want:     []Symbol{}, // nameフィールドがないため、テスト関数も出力されない
			wantErr:  false,
		},
		{
			name:     "大文字小文字を区別しない",
			filePath: "testdata/case_insensitive_test.go",
			want: []Symbol{
				{
					Name:   "TestCaseInsensitive",
					Detail: "test function",
					Kind:   SymbolKindFunction,
					Children: []Symbol{
						{
							Name:   "大文字NAME",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "混合Name",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "小文字name",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "空のファイルパス",
			filePath: "",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "存在しないファイル",
			filePath: "testdata/non_existent.go",
			want:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				t.Logf("Parse() error = %v", err)
			}

			if !tt.wantErr {
				if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(Symbol{}, "Range")); diff != "" {
					t.Errorf("Parse() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
