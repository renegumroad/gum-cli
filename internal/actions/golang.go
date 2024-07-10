package actions

import (
	"github.com/renegumroad/gum-cli/internal/cli/homebrew"
	"github.com/renegumroad/gum-cli/internal/systeminfo"
)

type GolangAction struct {
}

func NewGolangAction() *GolangAction {
	return &GolangAction{}
}

func (a *GolangAction) Name() string {
	return "golang"
}

func (a *GolangAction) Identifier() string {
	return "golang"
}

func (a *GolangAction) IsPublic() bool {
	return true
}

func (a *GolangAction) Platforms() []systeminfo.Platform {
	return []systeminfo.Platform{systeminfo.Darwin, systeminfo.Linux}
}

func (a *GolangAction) Deps() []Action {
	return []Action{
		NewBrewAction(
			[]homebrew.Package{
				{Name: "go"},
				{Name: "goreleaser"},
				{Name: "golangci-lint"},
				{Name: "go-task"},
				{Name: "mockery"},
				{Name: "gopls"},
			}),
	}
}

func (a *GolangAction) Validate() error {
	return nil
}

func (a *GolangAction) ShouldRun() bool {
	return depsShouldRun(a.Deps())
}

func (a *GolangAction) Run() error {
	return nil
}
