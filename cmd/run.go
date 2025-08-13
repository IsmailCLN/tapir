package cmd

import (
	"context"
	"fmt"

	"github.com/IsmailCLN/tapir/internal/parser"
	"github.com/IsmailCLN/tapir/internal/runner"
	"github.com/IsmailCLN/tapir/internal/ui"
	"github.com/spf13/cobra"
)

var file string

var runCmd = &cobra.Command{
	Use:   "run [suite.yaml]",
	Short: "Run a test suite",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {

		var path string
		switch {
		case file != "":
			path = file
		case len(args) == 1:
			path = args[0]
		default:
			return fmt.Errorf("no suite file given (use arg or -f)")
		}

		suites, err := parser.LoadTestSuite(path)
		if err != nil {
			return err
		}

		ctx := context.Background()
		results, err := runner.Run(ctx, suites)
		if err != nil {
			return err
		}

		return ui.Render([]string{path}, results)
	},
}
