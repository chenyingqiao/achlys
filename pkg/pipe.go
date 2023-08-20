package pkg

import (
	"io"
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

func (p *processPipe) AttachExtraFilesToProcess(cmd *exec.Cmd) {
	cmd.ExtraFiles = []*os.File{p.r}
}

// Send发送字符数组
func (p *processPipe) Send(cmd []string) {
	defer func() {
		p.cleanup()
		p.w.Close()
	}()
	p.w.WriteString(strings.Join(cmd, " "))
}

func (p *processPipe) cleanup() {
	p.r = nil
	p.w = nil
}

func ReadExtraFd() ([]string, error) {
	pipe := os.NewFile(uintptr(3), "pipe")
	bytes, err := io.ReadAll(pipe)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(bytes), " "), nil
}
