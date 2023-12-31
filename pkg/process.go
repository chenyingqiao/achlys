package pkg

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

const DefaultCloneFlags uintptr = syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC | syscall.CLONE_NEWUSER

type PorcessOption struct {
	KeepStedin bool
	TTY        bool
	Args       []string
	PipeArgs []string
	Cloneflags uintptr //进程克隆参数
}

type ProcessIface interface {
	Run(command string, opt PorcessOption)
}

type Process struct{}

func NewProcess() *Process {
	return &Process{}
}

// Run 在命名空间中执行一个新的进程
func (*Process) Run(command string, opt PorcessOption) error {
	log.Printf("start run command in container: %s %v\n", command, opt.Args)
	if opt.Cloneflags == 0 {
		log.Printf("start run command in container: %s %v\n", command, opt.Args)
		opt.Cloneflags = DefaultCloneFlags
	}
	cmd := exec.Command(command,opt.Args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: opt.Cloneflags,
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
	//通过额外的fd传入命令
	processPipe,err := NewProcessPipe()
	if err != nil {
		log.Printf("创建管道失败%s", err.Error())
		return err
	}
	err = processPipe.AttachExtraFilesToProcess(cmd)
	if err != nil {
		log.Printf("添加拓展管道文件失败%s\n", err.Error())
		return err
	}

	//这边要添加输出到文件的功能
	if opt.TTY {
		log.Println("启用tty")
		if opt.KeepStedin {
			cmd.Stdin = os.Stdin
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	log.Printf("进程开始执行")
	if err := cmd.Start(); err != nil {
		log.Printf("执行进程失败: %s", err.Error())
		return err
	}

	// 发送命令参数
	log.Printf("管道发送：%v",opt.PipeArgs)
	err = processPipe.Send(opt.PipeArgs)
	if err != nil {
		log.Printf("管道发送失败%s\n", err.Error())
		return err
	}
	log.Printf("管道完成")

	// 等待进程完成
	if err := cmd.Wait(); err != nil {
		log.Printf("进程等待失败: %s", err.Error())
		return err
	}

	log.Printf("进程已完成")
	return nil
}
