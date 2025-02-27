package actions

import (
	"github.com/pkg/errors"
	"github.com/renegumroad/gum-cli/internal/cli/homebrew"
	"github.com/renegumroad/gum-cli/internal/systeminfo"
)

type BrewAction struct {
	brew     homebrew.Client
	packages []homebrew.Package
}

func NewBrewAction(packages []homebrew.Package) Action {
	return newBrewActionWithClient(packages, homebrew.New())
}

func newBrewActionWithClient(packages []homebrew.Package, brew homebrew.Client) *BrewAction {
	return &BrewAction{
		brew:     brew,
		packages: packages,
	}
}

func (act *BrewAction) Name() string {
	return "brew"
}

func (act *BrewAction) Identifier() string {
	id := "brew"
	for _, pkg := range act.packages {
		id += "-" + pkg.Name
	}

	return id
}

func (act *BrewAction) Platforms() []systeminfo.Platform {
	return []systeminfo.Platform{systeminfo.Darwin, systeminfo.Linux}
}

func (act *BrewAction) Deps() []Action {
	return []Action{
		NewBrewEnsureAction(),
	}
}

func (act *BrewAction) Validate() error {
	missingNameCount := 0
	for _, pkg := range act.packages {
		if pkg.Name == "" {
			missingNameCount++
		}
	}

	var err error
	if missingNameCount > 0 {
		err = errors.Errorf("Failed %s action validation: %d package(s) missing name.", act.Name(), missingNameCount)
	} else if len(act.packages) == 0 {
		err = errors.Errorf("Failed %s action validation: no packages specified.", act.Name())
	}

	return err
}

func (act *BrewAction) IsPublic() bool {
	return true
}

func (act *BrewAction) ShouldRun() bool {
	if depsShouldRun(act.Deps()) {
		return true
	}

	for _, pkg := range act.packages {
		if !act.brew.IsInstalled(pkg) {
			return true
		}
	}

	return false
}

func (act *BrewAction) Run() error {
	for _, pkg := range act.packages {
		err := act.brew.EnsureInstalled(pkg)
		if err != nil {
			return err
		}
	}

	return nil
}
