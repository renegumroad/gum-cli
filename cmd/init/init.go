package init

import (
	"github.com/renehernandez/gum-cli/internal/commands"
	"github.com/renehernandez/gum-cli/internal/utils"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	impl := commands.NewInitCmd()

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize your local machine for gum usage.",
		Long: `Initialize your local machine for gum usage. It will make sure the following conditions are met:
* ~/.gum folder exists
* gum is added to the path at the corresponding location for your OS

This command is always safe to run, even if you have already initialized your machine.`,
		PreRun: func(_ *cobra.Command, _ []string) {
			utils.CheckFatalError(impl.Validate())
		},
		Run: func(_ *cobra.Command, _ []string) {
			utils.CheckFatalError(impl.Run())
		},
	}

	return cmd
}
