package sandbox

type Sandbox interface {
	Start() (int, error)
	Kill() error
}
