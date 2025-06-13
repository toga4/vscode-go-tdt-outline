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
			name:     "basic table-driven test",
			filePath: "testdata/basic_table_test.go",
			want: []Symbol{
				{
					Name:   "TestExample",
					Detail: "test function",
					Kind:   SymbolKindFunction,
					Children: []Symbol{
						{
							Name:   "normal case",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "zero value",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "multiple test functions",
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
			name:     "support for various field names",
			filePath: "testdata/various_fields_test.go",
			want: []Symbol{
				{
					Name:   "TestVariousFields",
					Detail: "test function",
					Kind:   SymbolKindFunction,
					Children: []Symbol{
						{
							Name:   "description field",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "title field",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "scenario field",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "testName field",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "multiple test tables",
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
			name:     "non-test functions are ignored",
			filePath: "testdata/non_test_functions.go",
			want:     []Symbol{},
			wantErr:  false,
		},
		{
			name:     "test cases without name field are ignored",
			filePath: "testdata/no_name_field_test.go",
			want:     []Symbol{}, // No test function output because there's no name field
			wantErr:  false,
		},
		{
			name:     "case insensitive field matching",
			filePath: "testdata/case_insensitive_test.go",
			want: []Symbol{
				{
					Name:   "TestCaseInsensitive",
					Detail: "test function",
					Kind:   SymbolKindFunction,
					Children: []Symbol{
						{
							Name:   "uppercase NAME",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "mixed case Name",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "lowercase name",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "empty file path",
			filePath: "",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "non-existent file",
			filePath: "testdata/non_existent.go",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "typed test cases",
			filePath: "testdata/typed_test_cases.go",
			want: []Symbol{
				{
					Name:   "TestTypedStruct",
					Detail: "test function",
					Kind:   SymbolKindFunction,
					Children: []Symbol{
						{
							Name:   "normal case: basic scenario",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "normal case: zero value scenario",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "error case: invalid input",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
					},
				},
				{
					Name:   "TestTypeAlias",
					Detail: "test function",
					Kind:   SymbolKindFunction,
					Children: []Symbol{
						{
							Name:   "type alias: case 1",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "type alias: case 2",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "type alias: case 3",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "map-based test cases",
			filePath: "testdata/map_test_cases.go",
			want: []Symbol{
				{
					Name:   "TestWithMap",
					Detail: "test function",
					Kind:   SymbolKindFunction,
					Children: []Symbol{
						{
							Name:   "normal case: basic scenario",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "normal case: zero value",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "error case: negative value",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
					},
				},
				{
					Name:   "TestSimpleMap",
					Detail: "test function",
					Kind:   SymbolKindFunction,
					Children: []Symbol{
						{
							Name:   "one",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "two",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "three",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
					},
				},
				{
					Name:   "TestTypedMap",
					Detail: "test function",
					Kind:   SymbolKindFunction,
					Children: []Symbol{
						{
							Name:   "empty string",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "hello world",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
						{
							Name:   "unicode",
							Detail: "test case",
							Kind:   SymbolKindStruct,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "non-go file extension",
			filePath: "testdata/basic_table_test.txt",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "file without extension",
			filePath: "testdata/basic_table_test",
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
