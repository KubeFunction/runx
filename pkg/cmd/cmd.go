package cmd

import (
	"github.com/spf13/cobra"

	"github.com/kubefunction/runx/pkg/cmd/templates"
)

func NewRunXCommand() *cobra.Command {
	cmds := &cobra.Command{
		Use:   "runx",
		Short: "runx desc",
		Long:  "runx desc",
		Run:   runHelp,
	}
	groups := templates.CommandGroups{
		{
			Message: "WASM Commands:",
			Commands: []*cobra.Command{
				newCmdWasm(),
			},
		},
		{
			Message: "Process Commands:",
			Commands: []*cobra.Command{
				newCmdProcess(),
			},
		},
	}
	groups.Add(cmds)
	return cmds
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
