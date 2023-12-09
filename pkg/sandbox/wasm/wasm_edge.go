package wasm

import (
	"k8s.io/klog"
)

type WasmEdgeSandboxConfig struct {
	WASMFile string
	Args     []string
	Pid      uint32
}
type WasmEdgeSandbox struct {
	Config *WasmEdgeSandboxConfig
}

func NewWasmEdgeSandbox(c *WasmEdgeSandboxConfig) *WasmEdgeSandbox {
	return &WasmEdgeSandbox{Config: c}
}
func (w *WasmEdgeSandbox) Start() (int, error) {
	klog.V(3).Infof("WasmEdge: wasm file %s and Args %s", w.Config.WASMFile, w.Config.Args)
	return 0, nil
}
func (w *WasmEdgeSandbox) Kill() error {
	return nil
}
func (w *WasmEdgeSandbox) CompleteConfig() {

}
