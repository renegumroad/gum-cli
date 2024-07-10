package actions

import (
	"github.com/pkg/errors"
	"github.com/renegumroad/gum-cli/internal/cli/homebrew"
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

func (act *BrewAction) Deps() []Action {
	return []Action{}
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

func (act *BrewAction) Public() bool {
	return true
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
