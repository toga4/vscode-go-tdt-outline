# Go TDD Outline Parser

A Go test file parser designed for VSCode extensions. It recognizes test cases from table-driven tests as individual symbols and displays them hierarchically in VSCode's outline view.

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

### Test Case Name Recognition

The parser automatically recognizes the following field names:
- `name`, `testName`, `desc`, `description`, `title`, `scenario`
- Case-insensitive comparison
- For map types, string keys are used as test case names

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

## Development & Testing

### Running Tests

```bash
# Run all tests
cd parser
go test ./...

# Run tests with coverage
go test -cover ./...
```

### Golden File Testing

This project uses Golden File Testing to track JSON output changes.

```bash
# Run golden file tests (compare with current output)
go test -run TestGoldenFiles

# Update golden files (when output format changes)
go test -run TestGoldenFiles -update
```

Golden files are stored in the `parser/testdata/golden/` directory.

### Usage Examples

```bash
# Verify parser operation
go run ./parser ./parser/internal/parser/testdata/map_test_cases.go | jq '.'

# Check part of the output
go run ./parser <file.go> | jq '.[0].children | map(.name)'
```
