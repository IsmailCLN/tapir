package cmd

import (
	"fmt"

	"github.com/IsmailCLN/tapir/internal/parser"
	"github.com/IsmailCLN/tapir/internal/ui"
	"github.com/spf13/cobra"
)

var file string

var runCmd = &cobra.Command{
	Use:   "run [suite.yaml]",
	Short: "Run test suites",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := file
		if path == "" {
			if len(args) == 0 {
				return fmt.Errorf("please provide a suite YAML path or use --file")
			}
			path = args[0]
		}

		if _, err := parser.LoadTestSuite(path); err != nil {
			return err
		}
		return ui.RenderStream([]string{path})
	},
}
