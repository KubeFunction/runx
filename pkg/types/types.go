package types

import (
	"github.com/kubefunction/runx/pkg/sandbox/libcontainer"
)

type ContainerInfo struct {
	libcontainer.ContainerState
	Labels      map[string]string `json:"labels,omitempty"`
	ContainerId string            `json:"container_id"`
}
