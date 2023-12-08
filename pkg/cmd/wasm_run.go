package cmd

import (
	"k8s.io/klog"
)

type WasmRunOption struct {
	WasmFile string
}

func (o *WasmRunOption) Run() error {
	klog.Infof("test cmd %s", o.WasmFile)
	return nil
}
