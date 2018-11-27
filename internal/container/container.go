package container

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// NewParentProc create parent init process for container
func NewParentProc(tty bool) (*exec.Cmd, *os.File, error) {
	// fork self to run init command
	cmd, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return nil, nil, err
	}

	r, w, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}

	initProc := exec.Command(cmd, "init")

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
		// user namespace
		// syscall.CLONE_NEWUSER,
	}

	// pass pipe read file to sub process
	initProc.ExtraFiles = append(initProc.ExtraFiles, r)

	// enable tty
	if tty {
		initProc.Stdin = os.Stdin
		initProc.Stdout = os.Stdout
		initProc.Stderr = os.Stderr
	}

	return initProc, w, nil
}

// RunContainerInitProc use syscall execev to takeover init process
func RunContainerInitProc() error {

	// get pipe
	pipe := os.NewFile(3, "pipe")

	// read commands
	bts, err := ioutil.ReadAll(pipe)
	if err != nil {
		return err
	}

	cmds := strings.Split(string(bts), " ")
	if len(cmds) == 0 {
		return errors.New("null command")
	}
	cmd := cmds[0]
	path, err := exec.LookPath(cmd)
	if err != nil {
		return err
	}

	// MS_NOEXEC: not run other proc, MS_NOSUID: not set uid, MS_NODEV: default
	mountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	// mount proc
	if err := syscall.Mount("proc", "/proc", "proc", uintptr(mountFlags), ""); err != nil {
		log.Println("[ERROR] mount err:", err)
	}

	// syscall.Exec will takeover init process
	if err := syscall.Exec(path, cmds, os.Environ()); err != nil {
		log.Println("err: ", err)
	}
	return nil
}
