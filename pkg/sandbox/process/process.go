package wasm

type ProcessSandbox struct {
}

func (w *ProcessSandbox) Start() (int, error) {
	return 0, nil
}
func (w *ProcessSandbox) Kill() error {
	return nil
}
