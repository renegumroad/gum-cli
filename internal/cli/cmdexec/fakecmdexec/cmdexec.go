package fakecmdexec

import (
	"github.com/renegumroad/gum-cli/internal/cli/cmdexec"
)

type SettableCommand interface {
	cmdexec.Command

	SetCmd(string)
	SetArgs([]string)
	SetEnv([]string)
}

type NoOpCommand struct {
	cmd  string
	args []string
	env  []string
}

func NewNoOpCommand() SettableCommand {
	return &NoOpCommand{
		cmd:  "",
		args: []string{},
		env:  []string{},
	}
}

func (c *NoOpCommand) Run() error {
	return nil
}

func (c *NoOpCommand) Stdout() string {
	return ""
}

func (c *NoOpCommand) Stderr() string {
	return ""
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

func NewCmdGenerator(command SettableCommand) cmdexec.CmdGenerator {
	return func(cmd string, args ...string) cmdexec.Command {
		command.SetCmd(cmd)
		command.SetArgs(args)
		return command
	}
}

func NewEnvCmdGenerator(command SettableCommand) cmdexec.EnvCmdGenerator {
	return func(cmd string, args, env []string) cmdexec.Command {
		command.SetCmd(cmd)
		command.SetArgs(args)
		command.SetEnv(env)
		return command
	}
}
