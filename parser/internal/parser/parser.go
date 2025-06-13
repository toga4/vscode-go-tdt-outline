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

	// Skip non-test functions
	// Valid: TestMyFunction, TestAPICall, Test_snake_case
	// Invalid: testMyFunction, MyTest, BenchmarkTest
	if !strings.HasPrefix(funcDecl.Name.String(), "Test") {
		return nil
	}

	// Skip functions with return values
	// Valid: func TestXxx(t *testing.T) {...}
	// Invalid: func TestXxx(t *testing.T) error {...}
	if funcDecl.Type.Results != nil {
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
	processedLiterals := make(map[*ast.CompositeLit]bool)

	// First pass: look for test table assignments
	// Pattern examples:
	//   tests := []struct{...}{...}     // anonymous struct slice
	//   tests := []Test{...}            // named type slice
	//   tests := Tests{...}             // type alias
	//   testCases := []*TestCase{...}   // pointer slice
	ast.Inspect(body, func(n ast.Node) bool {
		testCases := extractTestCasesFromAssignment(n, processedLiterals, fset)
		allTestCases = append(allTestCases, testCases...)
		return true
	})

	// Second pass: look for inline composite literals that weren't processed in assignments
	// Pattern examples:
	//   for _, tc := range []struct{...}{...} { ... }
	//   t.Run("group", func(t *testing.T) {
	//       for _, tc := range []Test{...} { ... }
	//   })
	ast.Inspect(body, func(n ast.Node) bool {
		compLit, ok := n.(*ast.CompositeLit)
		if !ok || processedLiterals[compLit] {
			return true
		}

		// Check if it's a slice or array type
		if _, isArray := compLit.Type.(*ast.ArrayType); !isArray {
			return true
		}

		testCases := extractTestCasesFromCompositeLit(compLit, fset)
		allTestCases = append(allTestCases, testCases...)
		return true
	})

	return allTestCases
}

// extractTestCasesFromAssignment extracts test cases from assignment statements
func extractTestCasesFromAssignment(n ast.Node, processedLiterals map[*ast.CompositeLit]bool, fset *token.FileSet) []Symbol {
	// Check if it's an assignment statement
	// Pattern: <variable> := <value>
	assign, ok := n.(*ast.AssignStmt)
	if !ok {
		return nil
	}

	// Validate assignment structure
	// Valid: tests := ...
	// Invalid: a, b := ..., tests, other := ...
	if len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
		return nil
	}

	// Check if the variable name suggests it's a test table
	ident, ok := assign.Lhs[0].(*ast.Ident)
	if !ok {
		return nil
	}

	varName := strings.ToLower(ident.Name)
	if !isTestTableVariableName(varName) {
		return nil
	}

	// Check if RHS is a composite literal
	// Valid patterns:
	//   []struct{...}{...}
	//   []Test{...}
	//   Tests{...} (where Tests is []Test)
	compLit, ok := assign.Rhs[0].(*ast.CompositeLit)
	if !ok {
		return nil
	}

	processedLiterals[compLit] = true
	return extractTestCasesFromCompositeLit(compLit, fset)
}

// isTestTableVariableName checks if a variable name suggests it contains test cases
func isTestTableVariableName(name string) bool {
	// Common patterns:
	//   tests, testCases, testcases
	//   cases, scenarios, examples
	//   tt (common abbreviation)
	//   tcs (test cases abbreviation)
	return strings.Contains(name, "test") ||
		strings.Contains(name, "case") ||
		strings.Contains(name, "scenario") ||
		strings.Contains(name, "example") ||
		name == "tt" ||
		name == "tcs"
}

// extractTestCasesFromCompositeLit extracts test cases from a composite literal
func extractTestCasesFromCompositeLit(compLit *ast.CompositeLit, fset *token.FileSet) []Symbol {
	testCases := []Symbol{}

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
		testCases = append(testCases, Symbol{
			Name:   testName,
			Detail: "test case",
			Kind:   SymbolKindStruct,
			Range:  toRange(startPos, endPos),
		})
	}

	return testCases
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
