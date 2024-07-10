package cmdexec

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/renegumroad/gum-cli/internal/log"
)

type Command struct {
	Cmd    string
	Args   []string
	Env    []string
	Stdout string
	Stderr string
}

func New(cmd string, args ...string) *Command {
	return NewWithEnv(cmd, args, []string{})
}

func NewWithEnv(cmd string, args, env []string) *Command {
	return &Command{
		Cmd:  cmd,
		Args: args,
		Env:  env,
	}
}

func (c *Command) Run() error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(c.Cmd, c.Args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, c.Env...)

	err := cmd.Run()
	c.Stdout = stdout.String()
	c.Stderr = stderr.String()

	if err != nil {
		log.Debugf("Cmd execution stdout: %s", c.Stdout)
		log.Debugf("Cmd execution stderr: %s", c.Stderr)
	}

	return err
}
