package cmdexec

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/renegumroad/gum-cli/internal/log"
)

type Command interface {
	Stdout() string
	Stderr() string
	Run() error
	Cmd() string
	Args() []string
	Env() []string
}

type command struct {
	cmd    string
	args   []string
	env    []string
	stdout string
	stderr string
}

func New(cmd string, args ...string) Command {
	return NewWithEnv(cmd, args, []string{})
}

func NewWithEnv(cmd string, args, env []string) Command {
	return &command{
		cmd:  cmd,
		args: args,
		env:  env,
	}
}

func (c *command) Run() error {
	log.Debugf("Running command: %s %v with env: %s", c.cmd, c.args, c.env)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(c.cmd, c.args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, c.env...)

	err := cmd.Run()
	c.stdout = stdout.String()
	c.stderr = stderr.String()

	if err != nil {
		log.Debugf("Cmd execution stdout: %s", c.stdout)
		log.Debugf("Cmd execution stderr: %s", c.stderr)
	}

	return err
}

func (c *command) Stdout() string {
	return c.stdout
}

func (c *command) Stderr() string {
	return c.stderr
}

func (c *command) Cmd() string {
	return c.cmd
}

func (c *command) Args() []string {
	return c.args
}

func (c *command) Env() []string {
	return c.env
}

type CmdGenerator func(cmd string, args ...string) Command

func NewCommandGenerator() CmdGenerator {
	return New
}

type EnvCmdGenerator func(cmd string, args, env []string) Command

func EnvCommandGenerator() EnvCmdGenerator {
	return NewWithEnv
}
