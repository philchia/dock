package subsystem

import (
	"bufio"
	"errors"
	"os"
	"path"
	"strings"
)

// ResourceConfig set resource limition
type ResourceConfig struct {
	MemoryLimit string
	CPUSet      string
	CPUShare    string
}

// Subsystem interface
type Subsystem interface {
	Name() string
	Set(path string, conf *ResourceConfig) error
	Apply(path string, pid int) error
	Remove(path string) error
}

// mux guard systems
var systems = map[string]Subsystem{}

// RegisterSubsystem add a subsystem
// NOTE: this function must only be called during initialization time (i.e. in
// an init() function)
func RegisterSubsystem(subsystem Subsystem) {
	systems[subsystem.Name()] = subsystem
}

// FindSubsystemMountPoint find hierarchy cgroup mount point for subsystem
func FindSubsystemMountPoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}

	return ""
}

// GetCgroupPath get cgroup path of subsystem
func GetCgroupPath(subsystem string, cgroup string, createIfNotExist bool) (string, error) {
	mountPoint := FindSubsystemMountPoint(subsystem)
	if mountPoint == "" {
		return "", errors.New("subsystem not found")
	}

	cgroupPath := path.Join(mountPoint, cgroup)
	if _, err := os.Stat(cgroupPath); err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}

		if !createIfNotExist {
			return "", errors.New("cgroup path not found")
		}

		if err := os.Mkdir(cgroupPath, 0755); err != nil {
			return "", err
		}
	}

	return cgroupPath, nil
}
