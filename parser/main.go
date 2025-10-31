// parser/main.go
package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/toga4/vscode-go-tdt-outline/parser/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <file_path|->", os.Args[0])
	}

	arg := os.Args[1]
	var symbols []parser.Symbol
	var err error

	if arg == "-" {
		// Read from stdin
		symbols, err = parser.Parse("<stdin>", os.Stdin)
	} else {
		// Read from file
		symbols, err = parser.ParseFile(arg)
	}
	if err != nil {
		log.Fatalf("Failed to parse: %v", err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(symbols); err != nil {
		log.Fatalf("Failed to encode symbols: %v", err)
	}
}
