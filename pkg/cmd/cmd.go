package cmd

import (
	"flag"

	"github.com/spf13/cobra"
	flags "github.com/spf13/pflag"
	"k8s.io/klog"

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
	cmds.Flags().SortFlags = false
	groups.Add(cmds)

	klog.InitFlags(nil)
	flag.Parse()
	flags.CommandLine.AddGoFlagSet(flag.CommandLine)
	return cmds
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
