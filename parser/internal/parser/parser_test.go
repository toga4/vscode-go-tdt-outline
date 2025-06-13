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
					Kind:   11,
					Children: []Symbol{
						{
							Name:   "正常系",
							Detail: "test case",
							Kind:   12,
						},
						{
							Name:   "ゼロ値",
							Detail: "test case",
							Kind:   12,
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
					Kind:   11,
					Children: []Symbol{
						{
							Name:   "test1",
							Detail: "test case",
							Kind:   12,
						},
						{
							Name:   "test2",
							Detail: "test case",
							Kind:   12,
						},
					},
				},
				{
					Name:   "TestSecond",
					Detail: "test function",
					Kind:   11,
					Children: []Symbol{
						{
							Name:   "test3",
							Detail: "test case",
							Kind:   12,
						},
						{
							Name:   "test4",
							Detail: "test case",
							Kind:   12,
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
					Kind:   11,
					Children: []Symbol{
						{
							Name:   "説明フィールド",
							Detail: "test case",
							Kind:   12,
						},
						{
							Name:   "タイトルフィールド",
							Detail: "test case",
							Kind:   12,
						},
						{
							Name:   "シナリオフィールド",
							Detail: "test case",
							Kind:   12,
						},
						{
							Name:   "テスト名フィールド",
							Detail: "test case",
							Kind:   12,
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
					Kind:   11,
					Children: []Symbol{
						{
							Name:   "table1-test1",
							Detail: "test case",
							Kind:   12,
						},
						{
							Name:   "table1-test2",
							Detail: "test case",
							Kind:   12,
						},
						{
							Name:   "table2-test1",
							Detail: "test case",
							Kind:   12,
						},
						{
							Name:   "table2-test2",
							Detail: "test case",
							Kind:   12,
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
					Kind:   11,
					Children: []Symbol{
						{
							Name:   "大文字NAME",
							Detail: "test case",
							Kind:   12,
						},
						{
							Name:   "混合Name",
							Detail: "test case",
							Kind:   12,
						},
						{
							Name:   "小文字name",
							Detail: "test case",
							Kind:   12,
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
		{
			name:     "型定義されたテストケース",
			filePath: "testdata/typed_test_cases.go",
			want: []Symbol{
				{
					Name:   "TestTypedStruct",
					Detail: "test function",
					Kind:   11,
					Children: []Symbol{
						{
							Name:   "正常系: 基本的なケース",
							Detail: "test case",
							Kind:   12,
						},
						{
							Name:   "正常系: ゼロ値のケース",
							Detail: "test case",
							Kind:   12,
						},
						{
							Name:   "異常系: 不正な入力",
							Detail: "test case",
							Kind:   12,
						},
					},
				},
				{
					Name:   "TestTypeAlias",
					Detail: "test function",
					Kind:   11,
					Children: []Symbol{
						{
							Name:   "型エイリアス: ケース1",
							Detail: "test case",
							Kind:   12,
						},
						{
							Name:   "型エイリアス: ケース2",
							Detail: "test case",
							Kind:   12,
						},
						{
							Name:   "型エイリアス: ケース3",
							Detail: "test case",
							Kind:   12,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "map型のテストケース",
			filePath: "testdata/map_test_cases.go",
			want: []Symbol{
				{
					Name:   "TestWithMap",
					Detail: "test function",
					Kind:   11,
					Children: []Symbol{
						{
							Name:   "正常系: 基本的なケース",
							Detail: "test case",
							Kind:   12,
						},
						{
							Name:   "正常系: ゼロ値",
							Detail: "test case",
							Kind:   12,
						},
						{
							Name:   "異常系: 負の値",
							Detail: "test case",
							Kind:   12,
						},
					},
				},
				{
					Name:   "TestSimpleMap",
					Detail: "test function",
					Kind:   11,
					Children: []Symbol{
						{
							Name:   "one",
							Detail: "test case",
							Kind:   12,
						},
						{
							Name:   "two",
							Detail: "test case",
							Kind:   12,
						},
						{
							Name:   "three",
							Detail: "test case",
							Kind:   12,
						},
					},
				},
				{
					Name:   "TestTypedMap",
					Detail: "test function",
					Kind:   11,
					Children: []Symbol{
						{
							Name:   "empty string",
							Detail: "test case",
							Kind:   12,
						},
						{
							Name:   "hello world",
							Detail: "test case",
							Kind:   12,
						},
						{
							Name:   "unicode",
							Detail: "test case",
							Kind:   12,
						},
					},
				},
			},
			wantErr: false,
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
