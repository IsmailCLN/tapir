package cmd

import (
	"fmt"

	"github.com/IsmailCLN/tapir/internal/httpclient"
	"github.com/IsmailCLN/tapir/internal/parser"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [yaml-file]",
	Short: "Run a YAML test suite",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		suitePath := args[0]

		suite, err := parser.LoadTestSuite(suitePath)
		if err != nil {
			fmt.Printf("Error loading test suite: %v\n", err)
			return
		}

		httpclient.RunAllTests(suite)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
