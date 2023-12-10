package sandbox

type WasmRuntimeType string

const (
	WasmEdgeRuntime WasmRuntimeType = "WasmEdge"
	WasmTimeRuntime WasmRuntimeType = "WasmTime"
)

type Sandbox interface {
	Init() (uint32, error)
	Start() (uint32, error)
	Kill() error
}
