package subsystem

import (
	"log"
)

// CgroupManager manage all subsystems
type CgroupManager struct {
	cgroup string
}

// NewCgroupManager create new CgroupManager
func NewCgroupManager(cgroup string) *CgroupManager {
	return &CgroupManager{
		cgroup: cgroup,
	}
}

// Set resource limit for subsystems
func (m *CgroupManager) Set(conf *ResourceConfig) error {
	for _, subsystem := range systems {
		if err := subsystem.Set(m.cgroup, conf); err != nil {
			return err
		}
	}

	return nil
}

// Apply add process's pid to all sub systems
func (m *CgroupManager) Apply(pid int) error {
	for _, subsystem := range systems {
		if err := subsystem.Apply(m.cgroup, pid); err != nil {
			return err
		}
	}

	return nil
}

// Destroy remove all subsystems
func (m *CgroupManager) Destroy() error {
	for _, subsystem := range systems {
		if err := subsystem.Remove(m.cgroup); err != nil {
			log.Println("[ERROR]:", err)
		}
	}

	return nil
}
