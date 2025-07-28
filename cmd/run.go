package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/IsmailCLN/tapir/internal/config"
	"github.com/IsmailCLN/tapir/internal/httpclient"
	"github.com/IsmailCLN/tapir/internal/parser"
	"github.com/spf13/cobra"
)

var (
	timeoutValue time.Duration
	failFast     bool
)

var runCmd = &cobra.Command{
	Use:   "run [yaml-file]",
	Short: "Run a YAML test suite",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config.HTTPTimeout = timeoutValue
		config.FailFast = failFast

		suitePath := args[0]
		suite, err := parser.LoadTestSuite(suitePath)
		if err != nil {
			fmt.Printf("Error loading test suite: %v\n", err)
			os.Exit(1)
		}

		httpclient.RunAllTests(suite)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().DurationVar(&timeoutValue, "timeout", 10*time.Second, "HTTP request timeout (e.g. 5s, 1500ms, 1m)")
	runCmd.Flags().BoolVar(&failFast, "fail-fast", false, "Stop executing further tests after the first failure")
}
