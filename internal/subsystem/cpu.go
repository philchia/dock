package subsystem

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

// static check
var _ Subsystem = (*cpuSubSystem)(nil)

func init() {
	RegisterSubsystem(new(cpuSubSystem))
}

type cpuSubSystem struct {
	path string
}

func (s *cpuSubSystem) Name() string {
	return "cpu"
}

func (s *cpuSubSystem) Set(cgroup string, res *ResourceConfig) error {
	cgroupPath, err := GetCgroupPath(s.Name(), cgroup, true)
	if err != nil {
		return err
	}

	s.path = cgroupPath

	if res.CPUShare != "" {
		cpuLimitPath := path.Join(cgroupPath, "cpu.shares")
		if err := ioutil.WriteFile(cpuLimitPath, []byte(res.CPUShare), 0644); err != nil {
			return err
		}
	}

	return nil
}

func (s *cpuSubSystem) Apply(cgroup string, pid int) error {

	tasksPath := path.Join(s.path, "tasks")
	return ioutil.WriteFile(tasksPath, []byte(strconv.Itoa(pid)), 0644)
}

func (s *cpuSubSystem) Remove(cgroup string) error {
	return os.RemoveAll(s.path)
}
