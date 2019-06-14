package container

import (
	"os"
	"os/exec"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

//fork 新进程，隔绝宿主环境
func NewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}

func Run(tty bool, command string) {
	parent := NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	parent.Wait()
	os.Exit(-1)
}

//代码执行到这里，说明容器的进程已经创建
//使用mount 挂载，以便后面通过ps等系统命令去查看进程资源
func RunContainerInitProcess(command string, args []string) error {
	log.Infof("command %s", command)

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}
