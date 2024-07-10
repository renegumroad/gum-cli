package actions

import "github.com/renegumroad/gum-cli/internal/systeminfo"

type BrewEnsureAction struct {
}

func NewBrewEnsureAction() *BrewEnsureAction {
	return &BrewEnsureAction{}
}

func (a *BrewEnsureAction) Name() string {
	return "brew_ensure"
}

func (a *BrewEnsureAction) Identifier() string {
	return "brew_ensure"
}

func (a *BrewEnsureAction) IsPublic() bool {
	return false
}

func (a *BrewEnsureAction) Platforms() []systeminfo.Platform {
	return []systeminfo.Platform{systeminfo.Darwin, systeminfo.Linux}
}

func (a *BrewEnsureAction) Deps() []Action {
	return []Action{
		NewXcodeAction(),
		NewScriptAction(&ScriptActionArgs{
			Title:   "Install Homebrew",
			Command: "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)",
			Test:    "brew --version",
		}),
	}
}

func (a *BrewEnsureAction) Validate() error {
	return nil
}

func (a *BrewEnsureAction) ShouldRun() bool {
	return depsShouldRun(a.Deps())
}

func (a *BrewEnsureAction) Run() error {
	return nil
}
