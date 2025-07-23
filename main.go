package main

import (
	"fmt"
	"os"
	"github.com/IsmailCLN/tapir/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tapir <path-to-yaml>")
		return
	}

	suite, err := utils.LoadTestSuite(os.Args[1])
	if err != nil {
		fmt.Printf("Error loading test suite: %v\n", err)
		return
	}

	for _, test := range suite.Tests {
		utils.RunTestCase(test)
	}
}
