package dev

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dev",
		Short:   "commands for development",
		Aliases: []string{"d", "development"},
	}

	cmd.AddCommand(newUpCmd())

	return cmd
}
