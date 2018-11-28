package container

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// NewParentProc create parent init process for container
func NewParentProc(tty bool, root string) (*exec.Cmd, *os.File, error) {
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
	initProc.Dir = root
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
			syscall.CLONE_NEWIPC |
			// user namespace
			syscall.CLONE_NEWUSER,
		// run container as root
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
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

	if err := setupMount(); err != nil {
		return err
	}

	// syscall.Exec will takeover init process
	if err := syscall.Exec(path, cmds, os.Environ()); err != nil {
		return err
	}

	return nil
}

func pivotRoot(root string) error {
	// to make new root and old root in different fs, mount new root to new root use bind
	if err := syscall.Mount(root, root, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return err
	}

	// create dir to store old root
	pivotDir := filepath.Join(root, ".pivot_root")

	if err := os.Mkdir(pivotDir, 0700); err != nil {
		return err
	}

	// switch to a new fs
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return err
	}

	// change dir to root
	if err := syscall.Chdir("/"); err != nil {
		return err
	}

	// unmount old root
	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return err
	}

	// rm old root
	return os.RemoveAll(pivotDir)
}

func mountProc(newroot string) error {
	target := filepath.Join(newroot, "/proc")

	if err := os.MkdirAll(target, 0700); err != nil {
		return err
	}

	if err := syscall.Mount("proc", target, "proc", 0, ""); err != nil {
		return err
	}

	return nil
}

func setupMount() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := mountProc(cwd); err != nil {
		return err
	}

	if err := pivotRoot(cwd); err != nil {
		return err
	}

	return nil
}
