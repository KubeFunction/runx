package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

func newCmdWasm() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wasm SUBCOMMAND",
		Short: "wasm run|kill",
		Long:  "wasm run|kill",
		Run:   runHelp,
	}
	cmd.AddCommand(newWasmRun())
	cmd.AddCommand(newWasmKill())
	return cmd
}

func newWasmRun() *cobra.Command {
	o := WasmRunOption{}
	cmd := &cobra.Command{
		Use:                   "run -f FILE",
		DisableFlagsInUseLine: true,
		Short:                 "Run wasm file with WasmEdge",
		Long:                  "Run wasm file with WasmEdge",
		Run: func(cmd *cobra.Command, args []string) {
			if err := o.Run(); err != nil {
				klog.Errorf("wasm run cmd error %v", err)
			}
		},
	}
	cmd.Flags().StringVarP(&o.WasmFile, "file", "f", o.WasmFile, "The path of WASM file")
	return cmd
}

func newWasmKill() *cobra.Command {
	o := WasmKillOption{}
	cmd := &cobra.Command{
		Use:                   "kill -p PID",
		DisableFlagsInUseLine: true,
		Short:                 "Kill the wasm process",
		Long:                  "Kill the wasm process",
		Run: func(cmd *cobra.Command, args []string) {
			if err := o.Run(); err != nil {
				klog.Errorf("kill wasm process error %v", err)
			}
		},
	}
	cmd.Flags().IntVarP(&o.Pid, "pid", "p", o.Pid, "The pid of WASM process")
	return cmd
}
