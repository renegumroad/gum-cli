package fakecmdexec

import (
	"github.com/renegumroad/gum-cli/internal/cli/cmdexec"
)

type SettableCommand interface {
	cmdexec.Command

	SetCmd(string)
	SetArgs([]string)
	SetEnv([]string)
	SetStdout(string)
	SetStderr(string)
	SetRunError(error)
}

type NoOpOutputs struct {
	Stdout string
	Stderr string
	Err    error
}

type NoOpCommand struct {
	cmd    string
	args   []string
	env    []string
	stdout string
	stderr string
	err    error
}

func NewNoOpCommand() SettableCommand {
	return &NoOpCommand{
		cmd:  "",
		args: []string{},
		env:  []string{},
		err:  nil,
	}
}

func NewNoOpCommandWithOutputs(outputs *NoOpOutputs) SettableCommand {
	cmd := NewNoOpCommand()

	cmd.SetStdout(outputs.Stdout)
	cmd.SetStderr(outputs.Stderr)
	cmd.SetRunError(outputs.Err)

	return cmd
}

func (c *NoOpCommand) Run() error {
	return c.err
}

func (c *NoOpCommand) Stdout() string {
	return c.stdout
}

func (c *NoOpCommand) Stderr() string {
	return c.stderr
}

func (c *NoOpCommand) Cmd() string {
	return c.cmd
}

func (c *NoOpCommand) Args() []string {
	return c.args
}

func (c *NoOpCommand) Env() []string {
	return c.env
}

func (c *NoOpCommand) SetCmd(cmd string) {
	c.cmd = cmd
}

func (c *NoOpCommand) SetArgs(args []string) {
	c.args = args
}

func (c *NoOpCommand) SetEnv(env []string) {
	c.env = env
}

func (c *NoOpCommand) SetStdout(stdout string) {
	c.stdout = stdout
}

func (c *NoOpCommand) SetStderr(stderr string) {
	c.stderr = stderr
}

func (c *NoOpCommand) SetRunError(err error) {
	c.err = err
}

func NewCmdGenerator(commands ...SettableCommand) cmdexec.CmdGenerator {
	index := 0
	return func(cmd string, args ...string) cmdexec.Command {
		if index >= len(commands) {
			panic("No more commands available")
		}
		command := commands[index]
		command.SetCmd(cmd)
		command.SetArgs(args)
		index++
		return command
	}
}

func NewEnvCmdGenerator(commands ...SettableCommand) cmdexec.EnvCmdGenerator {
	index := 0

	return func(cmd string, args, env []string) cmdexec.Command {
		if index >= len(commands) {
			panic("No more commands available")
		}
		command := commands[index]
		command.SetCmd(cmd)
		command.SetArgs(args)
		command.SetEnv(env)
		index++
		return command
	}
}
