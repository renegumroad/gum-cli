/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	initCmd "github.com/renegumroad/gum-cli/cmd/init"
	"github.com/renegumroad/gum-cli/internal/log"
	"github.com/renegumroad/gum-cli/internal/version"

	"github.com/spf13/cobra"
)

var rootFlags = struct {
	LogLevel string
}{
	LogLevel: "info",
}

func rootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "gum",
		Short:   "gum is a command line interface for Gumroad developers",
		Version: version.VERSION,
		Run: func(_ *cobra.Command, _ []string) {
		},
	}

	rootCmd.PersistentFlags().StringVar(&rootFlags.LogLevel, "log-level", "info", "set log level")

	rootCmd.AddCommand(initCmd.Cmd())

	return rootCmd
}

func init() {
	cobra.OnInitialize(entrypoint)
}

func entrypoint() {
	err := log.Initialize(rootFlags.LogLevel)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd().Execute()
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
}
