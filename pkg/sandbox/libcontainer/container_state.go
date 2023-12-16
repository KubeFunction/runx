package libcontainer

import (
	"time"

	"github.com/opencontainers/runtime-spec/specs-go"
)

type ContainerState struct {
	specs.State
	Cmd     string
	Args    []string
	Created time.Time
}
