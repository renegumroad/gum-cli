package dev

import (
	"github.com/renegumroad/gum-cli/internal/actions"
	"github.com/renegumroad/gum-cli/internal/filesystem"
	"github.com/renegumroad/gum-cli/internal/gumconfig"
	"github.com/renegumroad/gum-cli/internal/log"
)

type UpImpl struct {
	fs      filesystem.Client
	config  *gumconfig.GumConfig
	handler *actions.ActionHandler
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

	parsedActions := []actions.Action{}

	for _, up := range impl.config.Up {
		var action actions.Action

		if up.Action != "" {
			action = actions.Get(string(up.Action))
		} else {
			action = actions.NewBrewAction(up.Brew)
		}

		parsedActions = append(parsedActions, action)
	}

	impl.handler = actions.NewActionHandler(parsedActions)

	if err := impl.handler.Validate(); err != nil {
		return err
	}

	log.Infof("%d action(s) validated successfully", len(impl.handler.Actions))

	return nil
}

func (impl *UpImpl) Run() error {
	log.Debugf("Running up command")

	if err := impl.handler.Run(); err != nil {
		return err
	}

	return nil
}
