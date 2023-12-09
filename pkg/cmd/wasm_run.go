package cmd

import (
	"k8s.io/klog"

	"github.com/kubefunction/runx/pkg/sandbox"
	"github.com/kubefunction/runx/pkg/sandbox/wasm"
)

type WasmRunOption struct {
	WasmFile string
	Args     []string
	Runtime  WasmRuntimeType // wasm runtime.such as WasmEdge、WasmTime, etc.
	Sandbox  sandbox.Sandbox
}

func (o *WasmRunOption) Run() error {
	pid, err := o.Sandbox.Start()
	klog.Infof("test cmd %d %v %s", pid, err, o.Args)
	return err
}

func (o *WasmRunOption) Complete() {
	switch o.Runtime {
	case WasmEdgeRuntime:
		c := &wasm.WasmEdgeSandboxConfig{
			WASMFile: o.WasmFile,
			Args:     o.Args,
		}
		o.Sandbox = wasm.NewWasmEdgeSandbox(c)
	default:

	}
}
