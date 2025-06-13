package main_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

func TestGoldenFiles(t *testing.T) {
	t.Parallel()

	// Build the parser binary
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "parser")

	cmd := exec.Command("go", "build", "-o", binaryPath, "./main.go")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build parser: %v\nOutput: %s", err, output)
	}

	// Test cases for golden file testing
	testCases := []struct {
		name       string
		inputFile  string
		goldenFile string
	}{
		{
			name:       "basic table driven test",
			inputFile:  "internal/parser/testdata/basic_table_test.go",
			goldenFile: "testdata/golden/basic_table_test.json",
		},
		{
			name:       "multiple test functions",
			inputFile:  "internal/parser/testdata/multiple_functions_test.go",
			goldenFile: "testdata/golden/multiple_functions_test.json",
		},
		{
			name:       "various field names",
			inputFile:  "internal/parser/testdata/various_fields_test.go",
			goldenFile: "testdata/golden/various_fields_test.json",
		},
		{
			name:       "multiple test tables",
			inputFile:  "internal/parser/testdata/multiple_tables_test.go",
			goldenFile: "testdata/golden/multiple_tables_test.json",
		},
		{
			name:       "non test functions",
			inputFile:  "internal/parser/testdata/non_test_functions.go",
			goldenFile: "testdata/golden/non_test_functions.json",
		},
		{
			name:       "no name field",
			inputFile:  "internal/parser/testdata/no_name_field_test.go",
			goldenFile: "testdata/golden/no_name_field_test.json",
		},
		{
			name:       "case insensitive matching",
			inputFile:  "internal/parser/testdata/case_insensitive_test.go",
			goldenFile: "testdata/golden/case_insensitive_test.json",
		},
		{
			name:       "typed test cases",
			inputFile:  "internal/parser/testdata/typed_test_cases.go",
			goldenFile: "testdata/golden/typed_test_cases.json",
		},
		{
			name:       "map based test cases",
			inputFile:  "internal/parser/testdata/map_test_cases.go",
			goldenFile: "testdata/golden/map_test_cases.json",
		},
		{
			name:       "backtick strings",
			inputFile:  "internal/parser/testdata/backtick_strings_test.go",
			goldenFile: "testdata/golden/backtick_strings_test.json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Run parser on the input file
			cmd := exec.Command(binaryPath, tc.inputFile)
			output, err := cmd.Output()
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					t.Fatalf("Parser failed: %v\nStderr: %s", err, exitErr.Stderr)
				}
				t.Fatalf("Failed to run parser: %v", err)
			}

			// Format JSON for consistent comparison
			var formatted bytes.Buffer
			if err := json.Indent(&formatted, output, "", "  "); err != nil {
				t.Fatalf("Failed to format JSON: %v", err)
			}
			actual := formatted.Bytes()

			if *update {
				// Update golden file
				if err := os.MkdirAll(filepath.Dir(tc.goldenFile), 0755); err != nil {
					t.Fatalf("Failed to create golden file directory: %v", err)
				}
				if err := os.WriteFile(tc.goldenFile, actual, 0644); err != nil {
					t.Fatalf("Failed to update golden file: %v", err)
				}
				t.Logf("Updated golden file: %s", tc.goldenFile)
			} else {
				// Compare with golden file
				expected, err := os.ReadFile(tc.goldenFile)
				if err != nil {
					t.Fatalf("Failed to read golden file: %v", err)
				}

				if !bytes.Equal(expected, actual) {
					// Show diff for better debugging
					t.Errorf("Output mismatch for %s", tc.inputFile)
					t.Errorf("Diff:\n%s", diff(string(expected), string(actual)))
				}
			}
		})
	}
}

// diff returns a simple line-by-line diff of two strings
func diff(expected, actual string) string {
	expectedLines := strings.Split(expected, "\n")
	actualLines := strings.Split(actual, "\n")

	var result []string
	maxLines := max(len(actualLines), len(expectedLines))

	for i := range maxLines {
		var expectedLine, actualLine string
		if i < len(expectedLines) {
			expectedLine = expectedLines[i]
		}
		if i < len(actualLines) {
			actualLine = actualLines[i]
		}

		if expectedLine != actualLine {
			if expectedLine != "" {
				result = append(result, fmt.Sprintf("-%s", expectedLine))
			}
			if actualLine != "" {
				result = append(result, fmt.Sprintf("+%s", actualLine))
			}
		}
	}

	return strings.Join(result, "\n")
}
