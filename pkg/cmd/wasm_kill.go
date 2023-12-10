package cmd

import (
	"k8s.io/klog"

	"github.com/kubefunction/runx/pkg/sandbox"
	"github.com/kubefunction/runx/pkg/sandbox/wasm"
)

type WasmKillOption struct {
	wasmOptions
	Pid int
}

func (o *WasmKillOption) Run() error {
	return o.Sandbox.Kill()
}
func (o *WasmKillOption) Complete() {
	switch o.Runtime {
	case sandbox.WasmEdgeRuntime:
		c := &wasm.WasmEdgeSandboxConfig{
			Pid: o.Pid,
		}
		o.Sandbox = wasm.NewWasmEdgeSandbox(c)
	default:
		klog.Fatalf("not support the wasm runtime %s", o.Runtime)
	}
}
