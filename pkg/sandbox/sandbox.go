package sandbox

import (
	"fmt"

	"github.com/opencontainers/runtime-spec/specs-go"
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
	Kill() error
	List() ([]string, error)
	Sate() (*specs.State, error)
}
