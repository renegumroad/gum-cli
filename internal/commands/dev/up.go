package dev

import (
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
	return nil
}
