package update

import (
	"github.com/renegumroad/gum-cli/internal/commands/update"
	"github.com/renegumroad/gum-cli/internal/utils"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	impl := update.New()

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update gum version",
		Example: `  # Update gum binary
  gum update
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
