package wasm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/second-state/WasmEdge-go/wasmedge"
	"k8s.io/klog"

	"github.com/kubefunction/runx/pkg/sandbox"
	"github.com/kubefunction/runx/pkg/sandbox/libcontainer"
	"github.com/kubefunction/runx/pkg/sandbox/system"
	"github.com/kubefunction/runx/pkg/types"
)

type WasmEdgeSandboxConfig struct {
	WASMFile     string
	Args         []string
	Pid          int
	FunctionName string
	Detach       bool
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
	pid := cmd.Process.Pid
	if err = system.GenerateContainerRootPath(sandbox.WasmEdgeRuntimeRootPath, pid); err != nil {
		// try to close the process
		w.Config.Pid = pid
		_ = w.Kill()
		return 0, err
	}
	klog.V(3).Infof("WasmEdge: wasm process id %d", cmd.Process.Pid)
	if w.Config.Detach {
		return pid, nil
	}
	return 0, cmd.Wait()

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

	// write container info into config.json
	cmd, status, err := system.GetContainerCmdAndStatus(w.Config.Pid)
	containerInfo := &types.ContainerInfo{
		ContainerState: libcontainer.ContainerState{
			State: specs.State{
				Version:     "1.0", // todo
				Status:      status,
				Pid:         w.Config.Pid,
				ID:          strconv.Itoa(w.Config.Pid),
				Bundle:      w.Config.WASMFile,
				Annotations: nil,
			},
			Cmd: cmd,
		},
		ContainerId: strconv.Itoa(w.Config.Pid),
		Labels:      nil,
	}
	err = system.WriteContainerInfo(sandbox.WasmEdgeRuntimeRootPath, w.Config.Pid, containerInfo)
	return 0, err
}
func (w *WasmEdgeSandbox) Kill() error {
	// step 1. kill process by pid
	p, err := os.FindProcess(w.Config.Pid)
	if err != nil {
		return err
	}
	if err = p.Signal(syscall.SIGTERM); err != nil && !errors.Is(err, os.ErrProcessDone) {
		return err
	} else if errors.Is(err, os.ErrProcessDone) {
		klog.Warningf("stop wasm process error: %v", err)
	}
	// step 2. remove root dir of the runtime
	return os.RemoveAll(fmt.Sprintf("%s/%d", sandbox.WasmEdgeRuntimeRootPath, w.Config.Pid))
}

func (w *WasmEdgeSandbox) List() ([]string, error) {
	entries, err := ioutil.ReadDir(sandbox.WasmEdgeRuntimeRootPath)
	if err != nil {
		return nil, err
	}
	containers := make([]string, len(entries))
	for i, e := range entries {
		containers[i] = e.Name()
	}
	return containers, nil
}

func (w *WasmEdgeSandbox) Sate() (*libcontainer.ContainerState, error) {
	var status = specs.StateRunning
	state := specs.State{
		Version:     "1.0", // todo
		Status:      status,
		Pid:         w.Config.Pid,
		ID:          strconv.Itoa(w.Config.Pid),
		Bundle:      "",
		Annotations: nil,
	}
	containerSate := &libcontainer.ContainerState{
		State: state,
	}
	cmd, status, err := system.GetContainerCmdAndStatus(w.Config.Pid)
	containerSate.Cmd = cmd
	containerSate.State.Status = status
	return containerSate, err
}
