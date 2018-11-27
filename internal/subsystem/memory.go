package subsystem

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

// static check
var _ Subsystem = (*memorySubsystem)(nil)

func init() {
	RegisterSubsystem(new(memorySubsystem))
}

type memorySubsystem struct {
}

func (m *memorySubsystem) Name() string {
	return "memory"
}

func (m *memorySubsystem) Set(cgroup string, conf *ResourceConfig) error {
	cgroupPath, err := GetCgroupPath(m.Name(), cgroup, true)
	if err != nil {
		return err
	}

	if conf.MemoryLimit != "" {
		memoryLimitPath := path.Join(cgroupPath, "memory.limit_in_bytes")
		if err := ioutil.WriteFile(memoryLimitPath, []byte(conf.MemoryLimit), 0644); err != nil {
			return err
		}
	}

	return nil
}

func (m *memorySubsystem) Apply(cgroup string, pid int) error {
	cgroupPath, err := GetCgroupPath(m.Name(), cgroup, false)
	if err != nil {
		return err
	}

	tasksPath := path.Join(cgroupPath, "tasks")
	return ioutil.WriteFile(tasksPath, []byte(strconv.Itoa(pid)), 0644)
}

func (m *memorySubsystem) Remove(cgroup string) error {
	cgroupPath, err := GetCgroupPath(m.Name(), cgroup, false)
	if err != nil {
		return err
	}

	return os.Remove(cgroupPath)
}
