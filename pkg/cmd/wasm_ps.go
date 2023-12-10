package cmd

import (
	"k8s.io/klog"

	"github.com/kubefunction/runx/pkg/sandbox"
	"github.com/kubefunction/runx/pkg/sandbox/wasm"
)

type WasmPsOption struct {
	wasmOptions
}

func (o *WasmPsOption) Run() error {
	containers, err := o.Sandbox.List()
	if err != nil {
		return err
	}
	klog.Infof("wasm containers %s", containers)
	return nil
}
func (o *WasmPsOption) Complete() {
	switch o.Runtime {
	case sandbox.WasmEdgeRuntime:
		c := &wasm.WasmEdgeSandboxConfig{}
		o.Sandbox = wasm.NewWasmEdgeSandbox(c)
	default:
		klog.Fatalf("not support the wasm runtime %s", o.Runtime)
	}
}
