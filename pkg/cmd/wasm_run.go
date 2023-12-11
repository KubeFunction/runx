package cmd

import (
	"fmt"

	"k8s.io/klog"

	"github.com/kubefunction/runx/pkg/sandbox"
	"github.com/kubefunction/runx/pkg/sandbox/wasm"
)

type WasmRunOption struct {
	wasmOptions
	WasmFile string
	Args     []string
	Detach   bool
}

func (o *WasmRunOption) Run() error {
	pid, err := o.Sandbox.Init()
	if err != nil {
		klog.Errorf("run wasm file error %v", err)
		return err
	}
	fmt.Print(pid)
	return nil
}
func (o *WasmRunOption) RunWasm() error {
	_, err := o.Sandbox.Start()
	return err
}
func (o *WasmRunOption) Complete() {
	switch o.Runtime {
	case sandbox.WasmEdgeRuntime:
		c := &wasm.WasmEdgeSandboxConfig{
			WASMFile: o.WasmFile,
			Args:     o.Args,
			Detach:   o.Detach,
		}
		o.Sandbox = wasm.NewWasmEdgeSandbox(c)
	default:
		klog.Fatalf("not support the wasm runtime %s", o.Runtime)
	}
}
