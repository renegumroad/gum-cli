package dev

import (
	"github.com/renegumroad/gum-cli/internal/commands/dev"
	"github.com/renegumroad/gum-cli/internal/utils"
	"github.com/spf13/cobra"
)

func newUpCmd() *cobra.Command {
	impl := dev.NewUp()

	cmd := &cobra.Command{
		Use:   "up",
		Short: "configures your development dependencies.",
		Long: `Configures development dependencies declared in the gum.yml file in the current directory.

It will also make sure that some basic infrastructure components are already in place
    `,
		Example: `  # Configure the default dev environment for the project
  gum dev up
`,
		PreRun: func(_ *cobra.Command, _ []string) {
			utils.CheckFatalError(impl.Validate())
		},
		Run: func(_ *cobra.Command, _ []string) {
			utils.CheckFatalError(impl.Run())
		},
	}

	return cmd
}
