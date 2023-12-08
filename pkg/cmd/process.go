package cmd

import (
	"github.com/spf13/cobra"
)

func newCmdProcess() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "process SUBCOMMAND",
		Short: "not supported",
		Long:  "not supported",
		Run:   runHelp,
	}
	return cmd
}
