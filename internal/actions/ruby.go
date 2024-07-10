package actions

import (
	"github.com/renegumroad/gum-cli/internal/cli/homebrew"
	"github.com/renegumroad/gum-cli/internal/cli/rbenv"
	"github.com/renegumroad/gum-cli/internal/systeminfo"
)

type RubyAction struct {
}

func NewRubyAction() *RubyAction {
	return &RubyAction{}
}

func (a *RubyAction) Name() string {
	return "ruby"
}

func (a *RubyAction) Identifier() string {
	return "ruby"
}

func (a *RubyAction) IsPublic() bool {
	return true
}

func (a *RubyAction) Platforms() []systeminfo.Platform {
	return []systeminfo.Platform{systeminfo.Darwin, systeminfo.Linux}
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

func (a *RubyAction) ShouldRun() bool {
	return depsShouldRun(a.Deps())
}

func (a *RubyAction) Run() error {
	client := rbenv.New()

	return client.EnsureRubyInstalled()
}
