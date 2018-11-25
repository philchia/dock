package container

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

// NewParentProc create parent init process for container
func NewParentProc(tty bool, cmd string) *exec.Cmd {
	// fork self to run init command
	initProc := exec.Command("/proc/self/exe", "init", cmd)

	// setup namesapce ioslation
	initProc.SysProcAttr = &syscall.SysProcAttr{
		// host namespace
		Cloneflags: syscall.CLONE_NEWUTS |
			// pid namespace
			syscall.CLONE_NEWPID |
			// mount namespace
			syscall.CLONE_NEWNS |
			// net namespace
			syscall.CLONE_NEWNET |
			// ipc namespace
			syscall.CLONE_NEWIPC,
	}

	// enable tty
	if tty {
		initProc.Stdin = os.Stdin
		initProc.Stdout = os.Stdout
		initProc.Stderr = os.Stderr
	}

	return initProc
}

// RunContainerInitProc use syscall execev to takeover init process
func RunContainerInitProc(cmd string) error {

	// MS_NOEXEC: not run other proc, MS_NOSUID: not set uid, MS_NODEV: default
	mountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	// mount proc to enable ps to check pid
	syscall.Mount("proc", "/proc", "proc", uintptr(mountFlags), "")

	// syscall.Exec will takeover init process
	if err := syscall.Exec(cmd, []string{cmd}, os.Environ()); err != nil {
		log.Println("err: ", err)
		// return err
	}
	return nil
}
