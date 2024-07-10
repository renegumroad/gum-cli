package actions

import (
	"github.com/renegumroad/gum-cli/internal/cli/homebrew"
	"github.com/renegumroad/gum-cli/internal/cli/rbenv"
)

type RubyAction struct{}

func (a *RubyAction) Name() string {
	return "ruby"
}

func (a *RubyAction) Public() bool {
	return true
}

func (a *RubyAction) Deps() []Action {
	return []Action{
		NewBrewAction(
			[]homebrew.Package{
				{Name: "rbenv"},
			}),
	}
}

func (a *RubyAction) Validate() error {
	return nil
}

func (a *RubyAction) Run() error {
	client := rbenv.New()

	return client.EnsureRubyInstalled()
}
