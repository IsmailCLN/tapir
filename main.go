package main

import (
	"fmt"
	"os"

	"github.com/IsmailCLN/tapir/internal/httpclient"
	"github.com/IsmailCLN/tapir/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tapir <path-to-yaml>")
		return
	}

	suite, err := parser.LoadTestSuite(os.Args[1])
	if err != nil {
		fmt.Printf("Error loading test suite: %v\n", err)
		return
	}

	httpclient.RunAllTests(suite)
}
