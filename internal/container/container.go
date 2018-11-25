package container

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

// NewParentProc create parent init process for container
func NewParentProc(tty bool, cmd string) *exec.Cmd {
	initProc := exec.Command("/proc/self/exe", "init", cmd)

	// setup namesapce ioslation
	initProc.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET,
	}

	// run as root user in container
	initProc.SysProcAttr.Credential = &syscall.Credential{
		Uid: 1,
		Gid: 1,
	}

	if tty {
		initProc.Stdin = os.Stdin
		initProc.Stdout = os.Stdout
		initProc.Stderr = os.Stderr
	}

	return initProc
}

func RunContainerInitProc(cmd string) error {

	mountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(mountFlags), "")

	if err := syscall.Exec(cmd, []string{cmd}, os.Environ()); err != nil {
		log.Println("err: ", err)
		return err
	}
	return nil
}
