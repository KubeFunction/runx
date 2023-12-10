package wasm

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/second-state/WasmEdge-go/wasmedge"
	"k8s.io/klog"

	"github.com/kubefunction/runx/pkg/sandbox"
)

type WasmEdgeSandboxConfig struct {
	WASMFile     string
	Args         []string
	Pid          int
	FunctionName string
}
type WasmEdgeSandbox struct {
	Config *WasmEdgeSandboxConfig
}

func NewWasmEdgeSandbox(c *WasmEdgeSandboxConfig) *WasmEdgeSandbox {
	return &WasmEdgeSandbox{Config: c}
}
func (w *WasmEdgeSandbox) Init() (int, error) {
	//todo do something else

	klog.V(3).Infof("WasmEdge: wasm file %s and Args %s", w.Config.WASMFile, w.Config.Args)
	args := []string{"wasm", "run-wasm", "-f", w.Config.WASMFile, "-r", string(sandbox.WasmEdgeRuntime)}
	args = append(args, w.Config.Args...)
	runx, err := filepath.EvalSymlinks("/proc/self/exe")
	if err != nil {
		return 0, err
	}
	cmd := exec.Command(runx, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	klog.V(3).Infof("WasmEdge: exec wasm file process id %d", cmd.Process.Pid)
	if err := cmd.Wait(); err != nil {
		return 0, err
	}
	return cmd.Process.Pid, nil
}

func (w *WasmEdgeSandbox) Start() (int, error) {
	var conf = wasmedge.NewConfigure(wasmedge.REFERENCE_TYPES)
	conf.AddConfig(wasmedge.WASI)
	var vm = wasmedge.NewVMWithConfig(conf)
	var wasi = vm.GetImportModule(wasmedge.WASI)
	defer vm.Release()
	defer conf.Release()
	wasi.InitWasi(
		w.Config.Args,   // The args
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
