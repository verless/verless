package cli

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/verless/verless/core"
)

// newServeCmd creates the `verless serve` command.
func newServeCmd() *cobra.Command {
	var (
		options core.ServeOptions
	)

	buildCmd := cobra.Command{
		Use: "serve PROJECT",
		RunE: func(cmd *cobra.Command, args []string) error {
			errs := core.RunServe(args[0], options)

			if len(errs) == 1 {
				return errs[0]
			} else if len(errs) > 1 {
				return errors.Errorf("several errors occurred while serving: %v", errs)
			}

			return nil
		},
		Args: cobra.ExactArgs(1),
	}

	return &buildCmd
}
