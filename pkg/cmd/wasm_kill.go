package cmd

import (
	"k8s.io/klog"
)

type WasmKillOption struct {
	wasmOptions
	Pid int
}

func (o *WasmKillOption) Run() error {
	klog.Infof("test kill cmd %d", o.Pid)
	return nil
}
