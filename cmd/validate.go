package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/IsmailCLN/tapir/internal/adapter/yaml"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [yaml-file]",
	Short: "Validate a test YAML file for correct structure",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		suite, err := yaml.LoadTestSuite(path)
		if err != nil {
			fmt.Printf("❌ Failed to parse file: %v\n", err)
			os.Exit(1)
		}

		valid := true
		for i, t := range suite.Tests {
			if t.Name == "" {
				fmt.Printf("❌ Test %d is missing a 'name'\n", i+1)
				valid = false
			}
			if t.Method == "" {
				fmt.Printf("❌ Test '%s' is missing a 'method'\n", t.Name)
				valid = false
			} else if !isValidMethod(t.Method) { // ← eklenen kısım
				fmt.Printf("❌ Test '%s' has invalid HTTP method: (allowed: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS)\n", t.Method)
				valid = false
			}
			if t.URL == "" {
				fmt.Printf("❌ Test '%s' is missing a 'url'\n", t.Name)
				valid = false
			}
			if t.Expect.Status == 0 {
				fmt.Printf("❌ Test '%s' is missing expected 'status'\n", t.Name)
				valid = false
			} else if t.Expect.Status < 100 || t.Expect.Status > 599 {
				fmt.Printf("❌ Test '%d' has invalid status code (must be 100–599).\n", t.Expect.Status)
				valid = false
			}
		}

		if valid {
			fmt.Println("✅ YAML file is valid.")
		} else {
			fmt.Println("⚠️ YAML file has validation errors.")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func isValidMethod(m string) bool {
	switch strings.ToUpper(m) {
	case "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS":
		return true
	default:
		return false
	}
}
