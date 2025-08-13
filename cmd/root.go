package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "tapir",
		Short: "HTTP-test runner inspired by Great Expectations",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&initOut, "out", "o", "test-data/sample.yaml", "Output YAML path")
	initCmd.Flags().StringVarP(&initSuite, "name", "n", "sample", "Suite name")
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "Overwrite if file exists")
}
