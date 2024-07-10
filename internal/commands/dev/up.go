package dev

import (
	"github.com/renegumroad/gum-cli/internal/actions"
	"github.com/renegumroad/gum-cli/internal/filesystem"
	"github.com/renegumroad/gum-cli/internal/gumconfig"
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
	currentDir, err := impl.fs.CurrentDir()
	if err != nil {
		return err
	}

	impl.config, err = gumconfig.New(currentDir)
	if err != nil {
		return err
	}

	return impl.config.Validate()
}

func (impl *UpImpl) Run() error {
	var err error
	var action actions.Action

	for _, up := range impl.config.Up {
		if up.Action != "" {
			action = actions.Get(string(up.Action))
		} else if len(up.Brew) > 0 {
			action = actions.NewBrewAction(up.Brew)
		}

		err = action.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
