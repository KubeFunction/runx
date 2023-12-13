package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/kubefunction/runx/pkg/sandbox"
)

func newCmdWasm() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wasm SUBCOMMAND",
		Short: "wasm run|kill",
		Long:  "wasm run|kill",
		Run:   runHelp,
	}
	cmd.AddCommand(newWasmRun())
	cmd.AddCommand(newWasmRunDo())
	cmd.AddCommand(newWasmKill())
	cmd.AddCommand(newWasmPs())
	cmd.AddCommand(newWasmState())

	return cmd
}

func newWasmRun() *cobra.Command {
	o := &WasmRunOption{}
	cmd := &cobra.Command{
		Use:                   "run -f FILE [-d]",
		DisableFlagsInUseLine: true,
		Short:                 "Start wasm file with Wasm Runtime",
		Long:                  "Start wasm file with Wasm Runtime",
		Run: func(cmd *cobra.Command, args []string) {
			o.Args = args
			o.Complete()
			if err := o.Run(); err != nil {
				klog.Errorf("wasm run cmd error %v", err)
			}
		},
	}
	cmd.Flags().StringVarP(&o.WasmFile, "file", "f", "", "The path of WASM file")
	cmd.Flags().BoolVarP(&o.Detach, "detach", "d", false, "Run wasm process in background and print process ID")
	cmd.Flags().StringVarP((*string)(&o.Runtime), "runtime", "r", string(sandbox.WasmEdgeRuntime), "The wasm runtime.such as WasmEdge、WasmTime, etc.")
	return cmd
}
func newWasmRunDo() *cobra.Command {
	o := &WasmRunOption{}
	cmd := &cobra.Command{
		Use:                   "run-wasm -f FILE [-d]",
		DisableFlagsInUseLine: true,
		Short:                 "",
		Long:                  "",
		Run: func(cmd *cobra.Command, args []string) {
			o.Args = args
			o.Complete()
			if err := o.RunWasm(); err != nil {
				klog.Errorf("wasm run cmd error %v", err)
			}
		},
	}
	cmd.Flags().StringVarP(&o.WasmFile, "file", "f", "", "The path of WASM file")
	cmd.Flags().BoolVarP(&o.Detach, "detach", "d", false, "Run wasm process in background and print process ID")
	cmd.Flags().StringVarP((*string)(&o.Runtime), "runtime", "r", string(sandbox.WasmEdgeRuntime), "The wasm runtime.such as WasmEdge、WasmTime, etc.")
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
			o.Complete()
			if err := o.Run(); err != nil {
				klog.Errorf("kill wasm process error %v", err)
			}
		},
	}
	cmd.Flags().IntVarP(&o.Pid, "pid", "p", o.Pid, "The pid of WASM process")
	cmd.Flags().StringVarP((*string)(&o.Runtime), "runtime", "r", string(sandbox.WasmEdgeRuntime), "The wasm runtime.such as WasmEdge、WasmTime, etc.")
	return cmd
}

func newWasmPs() *cobra.Command {
	o := WasmPsOption{}
	cmd := &cobra.Command{
		Use:                   "ps ",
		DisableFlagsInUseLine: true,
		Short:                 "List all active wasm process",
		Long:                  "List all active wasm process",
		Run: func(cmd *cobra.Command, args []string) {
			o.Complete()
			if err := o.Run(); err != nil {
				klog.Errorf("list wasm process error %v", err)
			}
		},
	}
	cmd.Flags().StringVarP((*string)(&o.Runtime), "runtime", "r", string(sandbox.WasmEdgeRuntime), "The wasm runtime.such as WasmEdge、WasmTime, etc.")
	return cmd
}

func newWasmState() *cobra.Command {
	o := WasmStateOption{}
	cmd := &cobra.Command{
		Use:                   "state -p PID",
		DisableFlagsInUseLine: true,
		Short:                 "Get wasm process state by pid",
		Long:                  "Get wasm process state by pid",
		Run: func(cmd *cobra.Command, args []string) {
			o.Complete()
			if err := o.Run(); err != nil {
				klog.Errorf("get wasm process state error %v", err)
			}
		},
	}
	cmd.Flags().IntVarP(&o.Pid, "pid", "p", o.Pid, "The pid of WASM process")
	return cmd
}
