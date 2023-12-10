package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/kubefunction/runx/pkg/sandbox"
)

func newCmdWasm() *cobra.Command {
	o := wasmOptions{}
	cmd := &cobra.Command{
		Use:   "wasm SUBCOMMAND",
		Short: "wasm run|kill",
		Long:  "wasm run|kill",
		Run:   runHelp,
	}
	cmd.AddCommand(newWasmRun(o))
	cmd.AddCommand(newWasmRunDo(o))
	cmd.AddCommand(newWasmKill(o))
	cmd.AddCommand(newWasmPs(o))
	cmd.Flags().StringVarP((*string)(&o.Runtime), "runtime", "r", string(sandbox.WasmEdgeRuntime), "The wasm runtime.such as WasmEdge、WasmTime, etc.")
	return cmd
}

func newWasmRun(options wasmOptions) *cobra.Command {
	o := &WasmRunOption{}
	o.wasmOptions = options
	cmd := &cobra.Command{
		Use:                   "run -f FILE",
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
	cmd.Flags().StringVarP(&o.WasmFile, "file", "f", o.WasmFile, "The path of WASM file")
	return cmd
}
func newWasmRunDo(options wasmOptions) *cobra.Command {
	o := &WasmRunOption{}
	o.wasmOptions = options
	cmd := &cobra.Command{
		Use:                   "run-wasm -f FILE",
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
	cmd.Flags().StringVarP(&o.WasmFile, "file", "f", o.WasmFile, "The path of WASM file")
	return cmd
}
func newWasmKill(options wasmOptions) *cobra.Command {
	o := WasmKillOption{}
	o.wasmOptions = options
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

func newWasmPs(options wasmOptions) *cobra.Command {
	o := WasmPsOption{}
	o.wasmOptions = options
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
	return cmd
}
