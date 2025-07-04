# Go TDT Outline

Display Go table-driven test cases in VSCode's Outline view for better test navigation and organization.
## Features

Go TDT Outline enhances your Go testing experience by:

- **Automatic Test Case Detection**: Recognizes table-driven test patterns in your Go test files
- **Outline View Integration**: Displays individual test cases as navigable symbols in VSCode's Outline view
- **Hierarchical Organization**: Shows test functions and their test cases in a tree structure
- **Multiple Pattern Support**: Works with various test table patterns including:
  - Anonymous struct slices
  - Named struct slices
  - Type aliases
  - Map-based test cases
- **Smart Name Detection**: Automatically detects test case names from common field names (name, testName, desc, description, title, scenario)

### Screenshot

When you have a test file like this, the extension displays it in the Outline view.

<img width="782" alt="image" src="https://github.com/user-attachments/assets/734de4b5-ec18-49ee-97cd-12ea77318955" />


## Requirements

- Visual Studio Code v1.90.0 or higher
- Go files must have valid syntax
- No Go installation required (the extension includes a pre-built parser)

## Installation

### From VSCode Marketplace

1. Open VSCode
2. Go to Extensions (Ctrl+Shift+X / Cmd+Shift+X)
3. Search for "Go TDT Outline"
4. Click Install

### From VSIX file

1. Download the `.vsix` file from the [releases page](https://github.com/toga4/vscode-go-tdt-outline/releases)
2. In VSCode, go to Extensions → Views and More Actions (⋯) → Install from VSIX
3. Select the downloaded file

## Usage

The extension activates automatically when you open a Go file. Test cases will appear in the Outline view panel.

1. Open any Go test file containing table-driven tests
2. Open the Outline view (Explorer sidebar → Outline)
3. Navigate through your test cases by clicking on them

## Extension Settings

This extension contributes the following settings:

* `go-tdt-outline.timeout`: Timeout for parser execution in milliseconds (default: 10000)
* `go-tdt-outline.maxFileSize`: Maximum file size to analyze in bytes (default: 1048576 / 1MB)

## Development

### Prerequisites

- Node.js and pnpm
- Go
- Visual Studio Code

### Setup

```bash
# Clone the repository
git clone https://github.com/toga4/vscode-go-tdt-outline.git
cd vscode-go-tdt-outline

# Install dependencies
pnpm install

# Build the Go parser
pnpm run build-parser

# Compile TypeScript
pnpm run compile
```

### Testing

```bash
# Run all tests
pnpm test

# Run Go parser tests
cd parser && go test ./...

# Update test snapshots
pnpm run test:update-snapshots
```

### Building

```bash
# Package the extension
pnpm run package
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [vscode-extension API](https://code.visualstudio.com/api)
- Inspired by Go's testing best practices
