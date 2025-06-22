package main_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/sergi/go-diff/diffmatchpatch"
)

var update = flag.Bool("update", false, "update golden files")

// diff returns a line-by-line diff of two strings
func diff(expected, actual string) string {
	dmp := diffmatchpatch.New()
	a, b, c := dmp.DiffLinesToChars(actual, expected)
	diffs := dmp.DiffMain(a, b, false)
	diffs = dmp.DiffCharsToLines(diffs, c)

	var result []string
	for _, diff := range diffs {
		lines := strings.Split(diff.Text, "\n")
		for _, line := range lines[:len(lines)-1] {
			switch diff.Type {
			case diffmatchpatch.DiffEqual:
				result = append(result, "  "+line)
			case diffmatchpatch.DiffInsert:
				result = append(result, color.GreenString("+ "+line))
			case diffmatchpatch.DiffDelete:
				result = append(result, color.RedString("- "+line))
			}
		}
	}

	return strings.Join(result, "\n")
}

func TestMain(m *testing.M) {
	color.NoColor = false // force color output
	os.Exit(m.Run())
}

func TestGoldenFiles(t *testing.T) {
	t.Parallel()

	// Build the parser binary
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "parser")

	cmd := exec.Command("go", "build", "-o", binaryPath, "./main.go")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build parser: %v\nOutput: %s", err, output)
	}

	goldenDir := "testdata/golden"

	files, err := os.ReadDir("internal/parser/testdata")
	if err != nil {
		t.Fatalf("Failed to read testdata directory: %v", err)
	}
	for _, file := range files {
		inputFile := filepath.Join("internal/parser/testdata", file.Name())
		goldenFile := filepath.Join(goldenDir, file.Name()+".json")

		t.Run(inputFile, func(t *testing.T) {
			t.Parallel()

			// Run parser on the input file
			cmd := exec.Command(binaryPath, inputFile)
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
				if err := os.MkdirAll(filepath.Dir(goldenFile), 0755); err != nil {
					t.Fatalf("Failed to create golden file directory: %v", err)
				}
				if err := os.WriteFile(goldenFile, actual, 0644); err != nil {
					t.Fatalf("Failed to update golden file: %v", err)
				}
				t.Logf("Updated golden file: %s", goldenFile)
			} else {
				// Compare with golden file
				expected, err := os.ReadFile(goldenFile)
				if err != nil {
					t.Fatalf("Failed to read golden file: %v", err)
				}

				if !bytes.Equal(expected, actual) {
					// Show diff for better debugging
					t.Errorf("Output mismatch for %s", inputFile)
					t.Errorf("Diff:\n%s", diff(string(expected), string(actual)))
				}
			}
		})
	}
}
