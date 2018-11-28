package subsystem

import (
	"log"
)

// CgroupManager manage all subsystems
type CgroupManager struct {
	containerID string
}

// NewCgroupManager create new CgroupManager
func NewCgroupManager(containerID string) *CgroupManager {
	return &CgroupManager{
		containerID: containerID,
	}
}

// Set resource limit for subsystems
func (m *CgroupManager) Set(conf *ResourceConfig) error {
	for _, subsystem := range systems {
		if err := subsystem.Set(m.containerID, conf); err != nil {
			return err
		}
	}

	return nil
}

// Apply add process's pid to all sub systems
func (m *CgroupManager) Apply(pid int) error {
	for _, subsystem := range systems {
		if err := subsystem.Apply(m.containerID, pid); err != nil {
			return err
		}
	}

	return nil
}

// Destroy remove all subsystems
func (m *CgroupManager) Destroy() error {
	for _, subsystem := range systems {
		if err := subsystem.Remove(m.containerID); err != nil {
			log.Println("[ERROR]:", err)
		}
	}

	return nil
}
