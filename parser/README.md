# Go TDD Outline Parser

A parser that generates VS Code outline symbols from Go test files. It recognizes test cases from table-driven tests as individual symbols and displays them hierarchically in VS Code's outline view.

## Usage

```bash
# Execute the parser
go run ./parser <test_file.go>

# Display formatted JSON output
go run ./parser <test_file.go> | jq '.'
```

## Supported Test Patterns

### 1. Slice of Anonymous Structs
```go
tests := []struct {
    name string
    // ...
}{
    {name: "test1"},
    // ...
}
```

### 2. Slice of Named Structs
```go
type Test struct {
    name string
    // ...
}

data := []Test{
    {name: "test1"},
    // ...
}
```

### 3. Using Type Aliases
```go
type Tests []Test

scenarios := Tests{
    {name: "test1"},
    // ...
}
```

### 4. Map-based Test Cases
```go
// Map with struct values
tests := map[string]struct {
    input int
    want  int
}{
    "normal case": {input: 1, want: 2},
    "error case":  {input: -1, want: -1},
}

// Simple map
cases := map[string]int{
    "one": 1,
    "two": 2,
}
```

**Note**: 
- For slice types, there are no restrictions on variable names. All struct slices with a `name` field (or compatible fields) are detected.
- For map types, string keys are used as test case names.

## Implemented Improvements

### 1. Enhanced Error Handling
- Replaced `log.Fatal` with proper error returns
- Added file path validation (empty string and non-Go file checks)

### 2. Panic Prevention
- Changed string index access to `strings.HasPrefix`
- Improved implementation for safer execution

### 3. More Flexible Test Case Recognition
- Support for multiple field names: `name`, `testName`, `desc`, `description`, `title`, `scenario`
- Case-insensitive comparison using `strings.EqualFold`
- Support for multiple test tables
- Support for typed test cases

### 4. Improved Code Quality
- Added comprehensive comments
- Defined constants (`SymbolKindFunction`, `SymbolKindStruct`)
- Clear separation of function responsibilities

### 5. Added Tests
- Basic table-driven test parsing
- Multiple test functions and test tables
- Error case validation
- Support verification for various field names
- Validation of typed test cases

## Output Format

```json
[
  {
    "name": "TestExample",
    "detail": "test function",
    "kind": 11,
    "range": {...},
    "children": [
      {
        "name": "normal test case",
        "detail": "test case",
        "kind": 12,
        "range": {...}
      }
    ]
  }
]
```

## Future Enhancement Ideas

### 1. Support for Complex Type Definitions
Currently unsupported complex patterns:
```go
// Using types from external packages
tests := []somepackage.TestCase{
    {Name: "test1"},
}

// Complex type definitions with interfaces
type TestCase interface {
    Name() string
}
```

### 2. Subtest Support
Recognition of hierarchical structures from nested t.Run subtests

### 3. Performance Optimization
Improved processing speed for large test files

### 4. Customizable Configuration
- Configurable field names for recognition
- Output format selection
- Filtering functionality

## Development

### Running Tests

```bash
# Run all tests
cd parser
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific tests
go test -run TestParse ./internal/parser
```

### Golden File Testing

This project uses Golden File Testing to track JSON output changes.

```bash
# Run golden file tests (compare with current output)
cd parser
go test -run TestGoldenFiles

# Update golden files (when output format changes)
go test -run TestGoldenFiles -update

# Test specific cases only
go test -run "TestGoldenFiles/map_test_cases$"
```

Golden files are stored in the `parser/testdata/golden/` directory.

#### Why Use Golden File Testing?

1. **JSON Output Visualization**: Verify the exact JSON format expected by TypeScript extensions
2. **Change Tracking**: Clearly understand how code changes affect JSON output
3. **Backward Compatibility**: Prevent unintended output format changes

### Usage Examples

```bash
# Verify parser operation
go run ./parser ./parser/internal/parser/testdata/map_test_cases.go | jq '.'

# Check part of the output
go run ./parser <file.go> | jq '.[0].children | map(.name)'
``` 