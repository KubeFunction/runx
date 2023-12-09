package cmd

import (
	"k8s.io/klog"
)

type ProcessRunOption struct {
	Image string
	Mem   string
}

func (o *ProcessRunOption) Run() error {
	klog.Infof("test cmd %s %s", o.Image)
	return nil
}

func (o *ProcessRunOption) Complete(args []string) {
	o.Image = args[0]
}
