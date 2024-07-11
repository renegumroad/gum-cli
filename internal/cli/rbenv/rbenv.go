package rbenv

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/renegumroad/gum-cli/internal/cli/cmdexec"
	"github.com/renegumroad/gum-cli/internal/cli/homebrew"
	"github.com/renegumroad/gum-cli/internal/log"
)

type Client interface {
	IsRubyInstalled() bool
	EnsureRubyInstalled() error
}

type client struct {
	cmdGen cmdexec.CmdGenerator
	brew   homebrew.Client
}

func New() Client {
	return newClientWithComponents(
		cmdexec.NewCommandGenerator(),
		homebrew.New(),
	)
}

func newClientWithComponents(
	gen cmdexec.CmdGenerator,
	brew homebrew.Client,
) Client {
	return &client{
		cmdGen: gen,
		brew:   brew,
	}
}

func (c *client) EnsureRubyInstalled() error {
	if c.IsRubyInstalled() {
		log.Infof("Ruby version is already installed")
		return nil
	}

	if err := c.updateRubyBuild(); err != nil {
		return err
	}

	log.Infof("Installing ruby version")

	cmd := c.cmdGen("rbenv", "install", "--skip-existing")
	err := cmd.Run()

	if err != nil {
		return errors.Errorf("Failed ruby installation: %s", err)
	}

	return nil
}

func (c *client) IsRubyInstalled() bool {
	log.Debugln("Checking ruby version")

	cmd := c.cmdGen("rbenv", "version")
	err := cmd.Run()
	if err != nil {
		log.Debugf("Failed to check ruby version: %s", err)
		return false
	}

	return !strings.Contains(cmd.Stdout(), "not installed")
}

func (c *client) updateRubyBuild() error {
	log.Infof("Updating ruby-build")

	if err := c.brew.Upgrade(homebrew.Package{Name: "ruby-build"}); err != nil {
		return errors.Errorf("Failed to update ruby-build: %s", err)
	}

	return nil
}
