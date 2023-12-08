package sandbox

type Sandbox interface {
	Run() (int, error)
	Kill() error
}
