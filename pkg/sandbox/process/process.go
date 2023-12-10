package wasm

type ProcessSandbox struct {
}

func (w *ProcessSandbox) Start() (uint32, error) {
	return 0, nil
}
func (w *ProcessSandbox) Kill() error {
	return nil
}
func (w *ProcessSandbox) Init() (uint32, error) {
	return 0, nil
}
