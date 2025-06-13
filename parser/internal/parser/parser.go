package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"slices"
	"strconv"
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

// VS Code SymbolKind constants
const (
	SymbolKindFunction = 11 // VS Code's SymbolKind.Function
	SymbolKindStruct   = 12 // VS Code's SymbolKind.Struct
)

// Parse analyzes a Go file and extracts test functions with their test cases
func Parse(filePath string) ([]Symbol, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}
	if !strings.HasSuffix(filePath, ".go") {
		return nil, fmt.Errorf("file must have .go extension: %s", filePath)
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Go file %s: %w", filePath, err)
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
	var allTestCases []Symbol

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

		// Check if it's a map type
		if _, ok := compLit.Type.(*ast.MapType); ok {
			testCases := extractTestCasesFromMap(compLit, fset)
			allTestCases = append(allTestCases, testCases...)
			return true
		}

		// Otherwise, treat as slice/array
		testCases := extractTestCasesFromSlice(compLit, fset)
		allTestCases = append(allTestCases, testCases...)

		return true
	})

	return allTestCases
}

// extractTestCasesFromMap extracts test cases from map pattern
func extractTestCasesFromMap(compLit *ast.CompositeLit, fset *token.FileSet) []Symbol {
	var testCases []Symbol

	for _, elt := range compLit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		testName, ok := extractStringLiteral(kv.Key)
		if !ok {
			continue
		}

		testCases = append(testCases, createTestCaseSymbol(testName, kv, fset))
	}

	return testCases
}

// extractTestCasesFromSlice extracts test cases from slice/array pattern
func extractTestCasesFromSlice(compLit *ast.CompositeLit, fset *token.FileSet) []Symbol {
	var testCases []Symbol

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

		testCases = append(testCases, createTestCaseSymbol(testName, caseLit, fset))
	}

	return testCases
}

// createTestCaseSymbol creates a Symbol for a test case
func createTestCaseSymbol(testName string, node ast.Node, fset *token.FileSet) Symbol {
	startPos := fset.Position(node.Pos())
	endPos := fset.Position(node.End())
	return Symbol{
		Name:   testName,
		Detail: "test case",
		Kind:   SymbolKindStruct,
		Range:  toRange(startPos, endPos),
	}
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

		// Extract string literal value and remove quotes
		// Pattern: "test case name" -> test case name
		testName, ok := extractStringLiteral(kve.Value)
		if !ok {
			continue
		}
		return testName
	}

	return ""
}

// testNameFields contains field names commonly used for test case names
var testNameFields = []string{
	"name",
	"testName",
	"desc",
	"description",
	"title",
	"scenario",
}

// isTestNameField checks if a field name is commonly used for test case names
func isTestNameField(fieldName string) bool {
	return slices.ContainsFunc(testNameFields, func(name string) bool {
		return strings.EqualFold(fieldName, name)
	})
}

func extractStringLiteral(expr ast.Expr) (string, bool) {
	basicLit, ok := expr.(*ast.BasicLit)
	if !ok || basicLit.Kind != token.STRING {
		return "", false
	}

	unquoted, err := strconv.Unquote(basicLit.Value)
	if err != nil {
		return "", false
	}

	return unquoted, true
}

// toRange converts token positions to VS Code range format (0-indexed)
func toRange(start, end token.Position) Range {
	return Range{
		Start: Line{Line: start.Line - 1, Character: start.Column - 1},
		End:   Line{Line: end.Line - 1, Character: end.Column - 1},
	}
}
