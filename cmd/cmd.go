package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

const (
	DefaultRumTimeout = 3600
)

type Cmd struct {
	Cmd           string
	Timeout       time.Duration
	TerminateChan chan int
	Setpgid       bool
	command       *exec.Cmd
	stdout        bytes.Buffer
	stderr        bytes.Buffer
}

func RunCMD(command string) (string, error) {
	cmd, err := NewCmd(command)
	if err != nil {
		return "", err
	}
	if err = cmd.Run(); err != nil {
		return "", fmt.Errorf("local Command `%s` failed: %v", command, cmd.Stderr())
	}
	return cmd.Stdout(), nil
}

func NewCmd(command string) (*Cmd, error) {
	c := &Cmd{Cmd: command}
	if c.Timeout == 0*time.Second {
		c.Timeout = DefaultRumTimeout * time.Second
	}
	if c.TerminateChan == nil {
		c.TerminateChan = make(chan int)
	}
	cmd := exec.Command("/bin/bash", "-c", c.Cmd)
	if c.Setpgid {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}
	cmd.Stderr = &c.stderr
	cmd.Stdout = &c.stdout
	c.command = cmd
	return c, nil
}

func (c *Cmd) Run() error {
	if err := c.command.Start(); err != nil {
		return err
	}
	errChan := make(chan error)
	go func() {
		errChan <- c.command.Wait()
		defer close(errChan)
	}()
	var err error
	select {
	case err = <-errChan:
	case <-time.After(c.Timeout):
		err = c.terminate()
		if err == nil {
			err = fmt.Errorf("local Command run timeout, cmd `%s`, time`%v`", c.Cmd, c.Timeout)
		}
	case <-c.TerminateChan:
		err = c.terminate()
		if err == nil {
			err = fmt.Errorf("local Command is terminated, cmd `%s`", c.Cmd)
		}
	}
	return err
}

func (c *Cmd) Stderr() string {
	return strings.TrimSpace(c.stderr.String())
}

func (c *Cmd) Stdout() string {
	return strings.TrimSpace(c.stdout.String())
}

func (c *Cmd) terminate() error {
	if c.Setpgid {
		return syscall.Kill(-c.command.Process.Pid, syscall.SIGKILL)
	} else {
		return syscall.Kill(c.command.Process.Pid, syscall.SIGKILL)
	}
}
