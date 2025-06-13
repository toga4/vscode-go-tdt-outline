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
		log.Fatalf("Usage: %s <file_path>", os.Args[0])
	}

	filePath := os.Args[1]
	symbols, err := parser.Parse(filePath)
	if err != nil {
		log.Fatalf("Failed to parse file: %v", err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(symbols); err != nil {
		log.Fatalf("Failed to encode symbols: %v", err)
	}
}
