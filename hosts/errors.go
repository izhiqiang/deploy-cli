package hosts

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
)

type ErrRunCommand struct {
	GHost   *GHost
	Command string
	Err     error
	Status  int
}

func NewErrRunCommand(host *GHost, command string, err error) *ErrRunCommand {
	errCommand := &ErrRunCommand{
		GHost:   host,
		Command: command,
		Err:     err,
	}
	var exitErr *ssh.ExitError
	if errors.As(err, &exitErr) {
		errCommand.Status = exitErr.ExitStatus()
	}
	return errCommand
}

func (e *ErrRunCommand) Error() string {
	return fmt.Sprintf("failed to execute cmd: `%s` err: %#v", e.Command, e.Err)
}

type ErrGHostConnect struct {
	Host *GHost
	Err  error
}

func NewErrGHostConnect(host *GHost, err error) *ErrGHostConnect {
	return &ErrGHostConnect{Host: host, Err: err}
}
func (e *ErrGHostConnect) Error() string {
	if e.Host == nil {
		return "host uninitialized"
	}
	return fmt.Sprintf(
		"ssh connection host: %s, port: %d User: %s err: %s , Please check",
		e.Host.Host.Host,
		e.Host.Host.Port,
		e.Host.Host.User,
		e.Err.Error(),
	)
}

type ErrGHostConnects struct {
	Ghosts []GHost
	Errs   []error
}

func NewErrGHostConnects() *ErrGHostConnects {
	return &ErrGHostConnects{Ghosts: make([]GHost, 0), Errs: []error{}}
}

func (c *ErrGHostConnects) Add(ghost GHost, err error) {
	c.Ghosts = append(c.Ghosts, ghost)
	c.Errs = append(c.Errs, err)
}

func (c *ErrGHostConnects) Err() error {
	if len(c.Ghosts) <= 0 {
		return nil
	}
	return c
}

func (c *ErrGHostConnects) Error() string {
	var errStrings []string
	for i, ghost := range c.Ghosts {
		errStrings = append(errStrings, fmt.Sprintf("group: %s err : %s", ghost.Group, c.Errs[i].Error()))
	}
	return strings.Join(errStrings, "\n")
}
