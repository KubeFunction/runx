package wasm

import (
	"encoding/json"
	"os"

	"github.com/second-state/WasmEdge-go/wasmedge"
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
	var conf = wasmedge.NewConfigure(wasmedge.REFERENCE_TYPES)
	conf.AddConfig(wasmedge.WASI)
	var vm = wasmedge.NewVMWithConfig(conf)
	var wasi = vm.GetImportModule(wasmedge.WASI)
	defer vm.Release()
	defer conf.Release()
	wasi.InitWasi(
		os.Args[1:],     // The args
		os.Environ(),    // The envs
		[]string{".:."}, // The mapping directories
	)
	res, err := vm.RunWasmFile(w.Config.WASMFile, "_start")
	if err != nil {
		return 0, err
	}
	s, _ := json.Marshal(res)
	klog.V(3).Infof("wasm function output %s %d", s, wasi.WasiGetExitCode())
	return 0, nil
}
func (w *WasmEdgeSandbox) Kill() error {
	return nil
}
