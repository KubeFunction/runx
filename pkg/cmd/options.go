package cmd

import (
	"github.com/kubefunction/runx/pkg/sandbox"
)

type wasmOptions struct {
	Runtime sandbox.WasmRuntimeType // wasm runtime.such as WasmEdge、WasmTime, etc.
	Sandbox sandbox.Sandbox
}
