package xcode

import (
	"github.com/renegumroad/gum-cli/internal/cli/cmdexec"
	"github.com/renegumroad/gum-cli/internal/log"
)

type Client interface {
	IsInstalled() bool
	EnsureInstalled() error
}

type client struct {
	cmdGen cmdexec.CmdGenerator
}

func New() Client {
	return newClientWithComponents(
		cmdexec.NewCommandGenerator(),
	)
}

func newClientWithComponents(
	gen cmdexec.CmdGenerator,
) Client {
	return &client{
		cmdGen: gen,
	}
}

func (c *client) IsInstalled() bool {
	log.Debugln("Checking if xcode is installed")

	cmd := c.cmdGen("xcode-select", "-p")
	err := cmd.Run()
	return err == nil
}

func (c *client) EnsureInstalled() error {
	if c.IsInstalled() {
		log.Debugln("xcode is already installed")
		return nil
	}

	log.Debugln("Installing xcode")
	cmd := c.cmdGen("xcode-select", "--install")

	return cmd.Run()
}
