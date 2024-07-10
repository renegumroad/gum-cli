package actions

import "github.com/renegumroad/gum-cli/internal/homebrew"

type GolangAction struct {
}

func (a *GolangAction) Name() string {
	return "golang"
}

func (a *GolangAction) Public() bool {
	return true
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

func (a *GolangAction) Run() error {
	return nil
}
