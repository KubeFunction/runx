package main

import (
	"k8s.io/klog"

	"github.com/kubefunction/runx/pkg/cmd"
)

func main() {
	command := cmd.NewRunXCommand()
	if err := command.Execute(); err != nil {
		klog.Fatal(err)
	}
}
