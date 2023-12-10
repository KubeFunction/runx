package sandbox

import (
	"fmt"
)

type WasmRuntimeType string

const (
	WasmEdgeRuntime WasmRuntimeType = "WasmEdge"
	WasmTimeRuntime WasmRuntimeType = "WasmTime"
	RootPath        string          = "/var/lib/docker"
)

var (
	WasmEdgeRuntimeRootPath = fmt.Sprintf("%s/%s", RootPath, WasmTimeRuntime)
)

type Sandbox interface {
	Init() (int, error)
	Start() (int, error)
	Kill(pid int) error
	List() ([]string, error)
}
