package cmd

import (
	"k8s.io/klog"
)

type WasmKillOption struct {
	wasmOptions
	Pid int
}

func (o *WasmKillOption) Run() error {
	if err := o.Sandbox.Kill(); err != nil {
		klog.Errorf("remove dir error %v", err)
		return err
	}
	return nil
}
