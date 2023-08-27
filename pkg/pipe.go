package pkg

import (
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type processPipe struct {
	r, w *os.File
}

func NewProcessPipe() (*processPipe, error) {
	p := &processPipe{}
	var err error
	p.r, p.w, err = os.Pipe()
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *processPipe) AttachExtraFilesToProcess(cmd *exec.Cmd) error {
	if p.r == nil {
		return errors.New("读管道端为空")
	}
	cmd.ExtraFiles = []*os.File{p.r}
	return nil
}

// Send发送字符数组
func (p *processPipe) Send(cmd []string) error {
	if p.w == nil {
		return errors.New("写管道端为空")
	}
	defer func() {
		p.w.Close()
		p.cleanup()
	}()
	log.Printf("管道写入:%v", cmd)
	p.w.WriteString(strings.Join(cmd, " "))
	return nil
}

func (p *processPipe) cleanup() {
	p.r = nil
	p.w = nil
}

func ReadExtraFd() ([]string, error) {
	log.Printf("开始读取管道数据\n")
	pipe := os.NewFile(uintptr(3), "pipe")
	bytes, err := io.ReadAll(pipe)
	if err != nil {
		log.Printf("读取管道数据出错：%s\n", err.Error())
		return nil, err
	}
	return strings.Split(string(bytes), " "), nil
}
