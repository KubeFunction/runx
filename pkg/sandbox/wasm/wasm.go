package wasm

type WASMSandbox struct {
}

func (w *WASMSandbox) Run() (int, error) {
	return 0, nil
}
func (w *WASMSandbox) Kill() error {
	return nil
}
