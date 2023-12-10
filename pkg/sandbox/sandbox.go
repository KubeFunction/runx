package sandbox

import (
	"fmt"
)

type WasmRuntimeType string

const (
	WasmEdgeRuntime WasmRuntimeType = "WasmEdge"
	WasmTimeRuntime WasmRuntimeType = "WasmTime"
	RootPath        string          = "/var/lib/wasm"
)

var (
	WasmEdgeRuntimeRootPath = fmt.Sprintf("%s/%s", RootPath, WasmEdgeRuntime)
)

type Sandbox interface {
	Init() (int, error)
	Start() (int, error)
	Kill(pid int) error
	List() ([]string, error)
}
