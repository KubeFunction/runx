package sandbox

type WasmRuntimeType string

const (
	WasmEdgeRuntime WasmRuntimeType = "WasmEdge"
	WasmTimeRuntime WasmRuntimeType = "WasmTime"
)

type Sandbox interface {
	Init() (int, error)
	Start() (int, error)
	Kill() error
}
