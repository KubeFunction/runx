package cmd

import (
	"encoding/json"
	"fmt"

	"k8s.io/klog"

	"github.com/kubefunction/runx/pkg/sandbox"
	"github.com/kubefunction/runx/pkg/sandbox/wasm"
)

type WasmStateOption struct {
	wasmOptions
	Pid int
}

func (o *WasmStateOption) Run() error {
	state, err := o.Sandbox.Sate()
	if err != nil {
		klog.Warningf("get wasm process state error %v", err)
	}
	stateJson, err := json.Marshal(state)
	if err != nil {
		klog.Errorf("json marshal container state error %v", err)
		return err
	}
	fmt.Printf("%s", stateJson)
	return nil
}

func (o *WasmStateOption) Complete() {
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
