package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

func newCmdProcess() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "process SUBCOMMAND",
		Short: "not supported",
		Long:  "not supported",
		Run:   runHelp,
	}
	cmd.AddCommand(newProcessRun())
	return cmd
}
func newProcessRun() *cobra.Command {
	o := ProcessRunOption{}
	cmd := &cobra.Command{
		Use:                   "run IMAGE [-mem MemoryLimits]",
		DisableFlagsInUseLine: true,
		Short:                 "Start",
		Long:                  "Start",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				klog.Errorf("You must give an image ptah")
				cmd.Help()
				return
			}
			o.Complete(args)
			if err := o.Run(); err != nil {
				klog.Errorf("Process run cmd error %v", err)
			}
		},
	}
	cmd.Flags().StringVarP(&o.Mem, "mem", "m", o.Mem, "Memory limits")
	return cmd
}
