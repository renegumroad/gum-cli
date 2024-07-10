package dev

import (
	"github.com/renegumroad/gum-cli/internal/actions"
	"github.com/renegumroad/gum-cli/internal/filesystem"
	"github.com/renegumroad/gum-cli/internal/gumconfig"
	"github.com/renegumroad/gum-cli/internal/log"
)

type UpImpl struct {
	fs     filesystem.Client
	config *gumconfig.GumConfig
}

func NewUp() *UpImpl {
	return newUpWithComponents(filesystem.New())
}

func newUpWithComponents(fs filesystem.Client) *UpImpl {
	return &UpImpl{
		fs: fs,
	}
}

func (impl *UpImpl) Validate() error {
	log.Debugf("Validating up command")
	currentDir, err := impl.fs.CurrentDir()
	if err != nil {
		return err
	}

	impl.config, err = gumconfig.New(currentDir)
	if err != nil {
		return err
	}

	if err := impl.config.Validate(); err != nil {
		return err
	}

	var action actions.Action
	for _, up := range impl.config.Up {
		if up.Action != "" {
			action = actions.Get(string(up.Action))
		} else {
			action = actions.NewBrewAction(up.Brew)
		}

		if err := actions.Validate(action); err != nil {
			return err
		}
	}

	log.Infof("%d action(s) validated successfully", len(impl.config.Up))

	return nil
}

func (impl *UpImpl) Run() error {
	log.Debugf("Running up command")
	var action actions.Action

	for _, up := range impl.config.Up {
		if up.Action != "" {
			action = actions.Get(string(up.Action))
		} else if len(up.Brew) > 0 {
			action = actions.NewBrewAction(up.Brew)
		}

		if err := actions.Run(action); err != nil {
			return err
		}
	}

	return nil
}
