package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"slices"
	"strings"
)

// Symbol represents a code symbol in VS Code's outline format
type Symbol struct {
	Name     string   `json:"name"`
	Detail   string   `json:"detail"`
	Kind     int      `json:"kind"` // VS Code's SymbolKind enumeration
	Range    Range    `json:"range"`
	Children []Symbol `json:"children"`
}

// Range represents a text range in a file
type Range struct {
	Start Line `json:"start"`
	End   Line `json:"end"`
}

// Line represents a position in a file (0-indexed)
type Line struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// SymbolKind constants matching VS Code's SymbolKind enumeration
const (
	SymbolKindFunction = 11
	SymbolKindStruct   = 12
)

// testNameFields contains field names commonly used for test case names
var testNameFields = []string{
	"name",
	"testName",
	"desc",
	"description",
	"title",
	"scenario",
}

// Parse analyzes a Go file and extracts test functions with their test cases
func Parse(filePath string) ([]Symbol, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	symbols := []Symbol{}
	ast.Inspect(node, func(n ast.Node) bool {
		symbol := extractTestFunction(n, fset)
		if symbol != nil {
			symbols = append(symbols, *symbol)
			return false // Don't traverse into this function
		}
		return true
	})

	return symbols, nil
}

// extractTestFunction extracts a test function symbol if the node is a test function
func extractTestFunction(n ast.Node, fset *token.FileSet) *Symbol {
	// Check if node is a function declaration
	// Pattern: func TestXxx(t *testing.T) {...}
	funcDecl, ok := n.(*ast.FuncDecl)
	if !ok {
		return nil
	}

	// Skip non-test functions (requires name starts with "Test" and no return values)
	if !strings.HasPrefix(funcDecl.Name.String(), "Test") || funcDecl.Type.Results != nil {
		return nil
	}

	// Extract test cases from the function body
	testCases := extractTestCases(funcDecl.Body, fset)
	if len(testCases) == 0 {
		return nil
	}

	startPos := fset.Position(funcDecl.Pos())
	endPos := fset.Position(funcDecl.End())
	return &Symbol{
		Name:     funcDecl.Name.Name,
		Detail:   "test function",
		Kind:     SymbolKindFunction,
		Range:    toRange(startPos, endPos),
		Children: testCases,
	}
}

// extractTestCases finds and extracts test cases from a function body
func extractTestCases(body *ast.BlockStmt, fset *token.FileSet) []Symbol {
	allTestCases := []Symbol{}

	// Look for all composite literals
	// Pattern examples:
	//   tests := []struct{...}{...}              // slice literal
	//   tests := []Test{...}                     // slice of named type
	//   tests := Tests{...}                      // type alias (e.g., type Tests []Test)
	//   tests := map[string]struct{...}{...}     // map with string keys
	//   for _, tc := range []struct{...}{...}    // inline usage
	ast.Inspect(body, func(n ast.Node) bool {
		compLit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		// Check if it's a map type with string keys
		if mapType, ok := compLit.Type.(*ast.MapType); ok {
			// Check if key type is string
			if ident, ok := mapType.Key.(*ast.Ident); ok && ident.Name == "string" {
				// Process map entries
				for _, elt := range compLit.Elts {
					kv, ok := elt.(*ast.KeyValueExpr)
					if !ok {
						continue
					}

					// Extract test name from key (string literal)
					var testName string
					if basicLit, ok := kv.Key.(*ast.BasicLit); ok && basicLit.Kind == token.STRING {
						testName = strings.Trim(basicLit.Value, `"`)
					}
					if testName == "" {
						continue
					}

					// Get position from the key-value pair
					startPos := fset.Position(kv.Pos())
					endPos := fset.Position(kv.End())
					allTestCases = append(allTestCases, Symbol{
						Name:   testName,
						Detail: "test case",
						Kind:   SymbolKindStruct,
						Range:  toRange(startPos, endPos),
					})
				}
				return true
			}
		}

		// Process slice/array composite literals (existing logic)
		// Extract test cases from this composite literal
		// We check all composite literals since we can't always determine
		// if a type alias refers to a slice without type information
		for _, elt := range compLit.Elts {
			// Each element should be a struct literal
			// Pattern: {name: "test1", input: "value", want: "expected"}
			caseLit, ok := elt.(*ast.CompositeLit)
			if !ok {
				continue
			}

			testName := extractTestName(caseLit)
			if testName == "" {
				continue
			}

			startPos := fset.Position(caseLit.Pos())
			endPos := fset.Position(caseLit.End())
			allTestCases = append(allTestCases, Symbol{
				Name:   testName,
				Detail: "test case",
				Kind:   SymbolKindStruct,
				Range:  toRange(startPos, endPos),
			})
		}

		return true
	})

	return allTestCases
}

// extractTestName extracts the test name from a struct literal
func extractTestName(caseLit *ast.CompositeLit) string {
	for _, kv := range caseLit.Elts {
		// Skip non-key-value expressions
		// This handles both:
		//   {name: "test1", ...}     // key-value form
		//   {"test1", ...}           // positional form (skip)
		kve, ok := kv.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		// Check if the key is an identifier
		ident, ok := kve.Key.(*ast.Ident)
		if !ok {
			continue
		}

		// Check if the field name is one of the common test name fields
		// Examples: name, testName, desc, description, title, scenario
		if !isTestNameField(ident.Name) {
			continue
		}

		// Extract string literal value
		// Pattern: "test case name"
		basicLit, ok := kve.Value.(*ast.BasicLit)
		if !ok || basicLit.Kind != token.STRING {
			continue
		}

		// Remove quotes from string literal
		// "test name" -> test name
		return strings.Trim(basicLit.Value, `"`)
	}

	return ""
}

// isTestNameField checks if a field name is commonly used for test case names
func isTestNameField(fieldName string) bool {
	return slices.ContainsFunc(testNameFields, func(name string) bool {
		return strings.EqualFold(fieldName, name)
	})
}

// toRange converts token positions to VS Code range format (0-indexed)
func toRange(start, end token.Position) Range {
	return Range{
		Start: Line{Line: start.Line - 1, Character: start.Column - 1},
		End:   Line{Line: end.Line - 1, Character: end.Column - 1},
	}
}
