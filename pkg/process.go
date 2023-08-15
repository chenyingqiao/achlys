package pkg

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

const DefaultCloneFlags uintptr = syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC

type PorcessOption struct{
	KeepStedin bool
	TTY bool
	Args []string
	Cloneflags uintptr //进程克隆参数
}

type ProcessIface interface{
	Run(command string,opt PorcessOption)
}

type Process struct{}

func NewProcess() *Process {
	return &Process{}
}

//Run 在命名空间中执行一个新的进程
func (*Process) Run(command string,opt PorcessOption) error {
	if opt.Cloneflags == 0 {
		opt.Cloneflags = DefaultCloneFlags
	}
	cmd := exec.Command(command, opt.Args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: opt.Cloneflags,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID: 0,
				Size: 1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID: 0,
				Size: 1,
			},
		},
	}
	//这边要添加输出到文件的功能
	if opt.TTY {
		if opt.KeepStedin {
			cmd.Stdin = os.Stdin
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Start(); err != nil {
		log.Printf("执行进程失败: %s",err.Error())
		return err
	}
	cmd.Wait()
	os.Exit(-1)
	return nil
}

