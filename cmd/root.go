package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
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
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "file", "f",
		"test.yaml", "YAML file containing suites")
}
