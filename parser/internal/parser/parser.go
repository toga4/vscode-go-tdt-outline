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
		// Look for function declarations starting with "Test"
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		// Skip non-test functions
		if !strings.HasPrefix(funcDecl.Name.String(), "Test") {
			return true
		}

		// Skip functions with return values (test functions should not return anything)
		if funcDecl.Type.Results != nil {
			return true
		}

		// Extract test cases from the function body
		testCases := extractTestCases(funcDecl.Body, fset)

		if len(testCases) > 0 {
			startPos := fset.Position(funcDecl.Pos())
			endPos := fset.Position(funcDecl.End())
			// Add function symbol with test cases as children
			symbols = append(symbols, Symbol{
				Name:     funcDecl.Name.Name,
				Detail:   "test function",
				Kind:     SymbolKindFunction,
				Range:    toRange(startPos, endPos),
				Children: testCases,
			})
		}

		return false
	})

	return symbols, nil
}

// extractTestCases finds and extracts test cases from a function body
func extractTestCases(body *ast.BlockStmt, fset *token.FileSet) []Symbol {
	allTestCases := []Symbol{}

	ast.Inspect(body, func(n ast.Node) bool {
		// Look for slice literals (e.g., tests := []struct{...}{...})
		compLit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		// Check if it's a slice or array type
		_, isArray := compLit.Type.(*ast.ArrayType)
		if !isArray {
			return true
		}

		// Process each test case (struct literal)
		testCases := extractTestCasesFromCompositeLit(compLit, fset)
		allTestCases = append(allTestCases, testCases...)

		// Continue searching for more test tables
		return true
	})

	return allTestCases
}

// extractTestCasesFromCompositeLit extracts test cases from a composite literal
func extractTestCasesFromCompositeLit(compLit *ast.CompositeLit, fset *token.FileSet) []Symbol {
	testCases := []Symbol{}

	for _, elt := range compLit.Elts {
		caseLit, ok := elt.(*ast.CompositeLit)
		if !ok {
			continue
		}

		testName := extractTestName(caseLit)
		if testName != "" {
			startPos := fset.Position(caseLit.Pos())
			endPos := fset.Position(caseLit.End())
			testCases = append(testCases, Symbol{
				Name:   testName,
				Detail: "test case",
				Kind:   SymbolKindStruct,
				Range:  toRange(startPos, endPos),
			})
		}
	}

	return testCases
}

// extractTestName extracts the test name from a struct literal
func extractTestName(caseLit *ast.CompositeLit) string {
	for _, kv := range caseLit.Elts {
		kve, ok := kv.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		ident, ok := kve.Key.(*ast.Ident)
		if !ok {
			continue
		}

		// Check if the field name is one of the common test name fields
		if slices.ContainsFunc(testNameFields, func(fieldName string) bool { return strings.EqualFold(ident.Name, fieldName) }) {
			if basicLit, ok := kve.Value.(*ast.BasicLit); ok && basicLit.Kind == token.STRING {
				return strings.Trim(basicLit.Value, `"`)
			}
		}
	}

	return ""
}

// toRange converts token positions to VS Code range format (0-indexed)
func toRange(start, end token.Position) Range {
	return Range{
		Start: Line{Line: start.Line - 1, Character: start.Column - 1},
		End:   Line{Line: end.Line - 1, Character: end.Column - 1},
	}
}
