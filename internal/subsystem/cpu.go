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
}

func (s *cpuSubSystem) Name() string {
	return "cpu"
}

func (s *cpuSubSystem) Set(cgroup string, res *ResourceConfig) error {
	cgroupPath, err := GetCgroupPath(s.Name(), cgroup, true)
	if err != nil {
		return err
	}

	if res.CPUShare != "" {
		if err := ioutil.WriteFile(path.Join(cgroupPath, "cpu.shares"), []byte(res.CPUShare), 0644); err != nil {
			return err
		}
	}

	return nil
}

func (s *cpuSubSystem) Remove(cgroup string) error {
	cgroupPath, err := GetCgroupPath(s.Name(), cgroup, false)
	if err != nil {
		return err
	}

	return os.RemoveAll(cgroupPath)
}

func (s *cpuSubSystem) Apply(cgroup string, pid int) error {
	cgroupPath, err := GetCgroupPath(s.Name(), cgroup, false)
	if err != nil {
		return err

	}

	if err := ioutil.WriteFile(path.Join(cgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
		return err
	}
	return nil
}
