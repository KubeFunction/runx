package system

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/opencontainers/runtime-spec/specs-go"

	"github.com/kubefunction/runx/pkg/types"
)

// State is the status of a process.
type State rune

const ( // Only values for Linux 3.14 and later are listed here
	Dead        State = 'X'
	DiskSleep   State = 'D'
	Running     State = 'R'
	Sleeping    State = 'S'
	Stopped     State = 'T'
	TracingStop State = 't'
	Zombie      State = 'Z'
)

// String forms of the state from proc(5)'s documentation for
// /proc/[pid]/status' "State" field.
func (s State) String() string {
	switch s {
	case Dead:
		return "dead"
	case DiskSleep:
		return "disk sleep"
	case Running:
		return "running"
	case Sleeping:
		return "sleeping"
	case Stopped:
		return "stopped"
	case TracingStop:
		return "tracing stop"
	case Zombie:
		return "zombie"
	default:
		return fmt.Sprintf("unknown (%c)", s)
	}
}

// Stat_t represents the information from /proc/[pid]/stat, as
// described in proc(5) with names based on the /proc/[pid]/status
// fields.
type Stat_t struct {
	// PID is the process ID.
	PID uint

	// Name is the command run by the process.
	Name string

	// State is the state of the process.
	State State

	// StartTime is the number of clock ticks after system boot (since
	// Linux 2.6).
	StartTime uint64
}

// Stat returns a Stat_t instance for the specified process.
func Stat(pid int) (stat Stat_t, err error) {
	bytes, err := ioutil.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "stat"))
	if err != nil {
		return stat, err
	}
	return parseStat(string(bytes))
}

// GetProcessStartTime is deprecated.  Use Stat(pid) and
// Stat_t.StartTime instead.
func GetProcessStartTime(pid int) (string, error) {
	stat, err := Stat(pid)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", stat.StartTime), nil
}

func parseStat(data string) (stat Stat_t, err error) {
	// From proc(5), field 2 could contain space and is inside `(` and `)`.
	// The following is an example:
	// 89653 (gunicorn: maste) S 89630 89653 89653 0 -1 4194560 29689 28896 0 3 146 32 76 19 20 0 1 0 2971844 52965376 3920 18446744073709551615 1 1 0 0 0 0 0 16781312 137447943 0 0 0 17 1 0 0 0 0 0 0 0 0 0 0 0 0 0
	i := strings.LastIndex(data, ")")
	if i <= 2 || i >= len(data)-1 {
		return stat, fmt.Errorf("invalid stat data: %q", data)
	}

	parts := strings.SplitN(data[:i], "(", 2)
	if len(parts) != 2 {
		return stat, fmt.Errorf("invalid stat data: %q", data)
	}

	stat.Name = parts[1]
	_, err = fmt.Sscanf(parts[0], "%d", &stat.PID)
	if err != nil {
		return stat, err
	}

	// parts indexes should be offset by 3 from the field number given
	// proc(5), because parts is zero-indexed and we've removed fields
	// one (PID) and two (Name) in the paren-split.
	parts = strings.Split(data[i+2:], " ")
	var state int
	fmt.Sscanf(parts[3-3], "%c", &state)
	stat.State = State(state)
	fmt.Sscanf(parts[22-3], "%d", &stat.StartTime)
	return stat, nil
}

func WriteContainerInfo(runtimeRoot string, pid int, info *types.ContainerInfo) error {
	containerInfo, err := json.Marshal(info)
	if err != nil {
		return err
	}
	infoPath := fmt.Sprintf("%s/%d/config.json", runtimeRoot, pid)
	f, err := os.OpenFile(infoPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.WriteString(string(containerInfo))
	return err
}

func GenerateContainerRootPath(runtimeRoot string, pid int) error {
	return os.MkdirAll(fmt.Sprintf("%s/%d", runtimeRoot, pid), 0755)
}

func GetContainerCmdAndStatus(pid int) (string, specs.ContainerState, error) {
	var cmd string
	var status = specs.StateRunning
	cmd, err := os.Readlink("/proc/" + strconv.Itoa(pid) + "/exe")
	if err != nil {
		return "", specs.StateStopped, err
	}
	stat, err := Stat(pid)
	if err != nil {
		status = specs.StateStopped
	} else if stat.State == Zombie || stat.State == Dead {
		status = specs.StateStopped
	}
	return cmd, status, err
}
