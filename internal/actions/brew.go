package actions

import (
	"github.com/pkg/errors"
	"github.com/renegumroad/gum-cli/internal/homebrew"
)

type BrewAction struct {
	brew     homebrew.Client
	packages []homebrew.Package
	helper   *actionHelper
}

func NewBrewAction(packages []homebrew.Package) Action {
	return newBrewActionWithClient(packages, homebrew.New())
}

func newBrewActionWithClient(packages []homebrew.Package, brew homebrew.Client) *BrewAction {
	return &BrewAction{
		brew:     brew,
		packages: packages,
		helper:   &actionHelper{},
	}
}

func (act *BrewAction) Name() string {
	return "brew"
}

func (act *BrewAction) Deps() []Action {
	return []Action{}
}

func (act *BrewAction) Validate() error {
	act.helper.logValidationMsg(act)
	depsError := act.helper.validateDependencies(act.Deps())

	missingNameCount := 0
	for _, pkg := range act.packages {
		if pkg.Name == "" {
			missingNameCount++
		}
	}

	var err error
	if missingNameCount > 0 {
		err = errors.Errorf("Failed %s action validation: %d package(s) missing name.", act.Name(), missingNameCount)
	}

	if depsError != nil {
		err = errors.Errorf("%s\n%s", err, depsError)
	}

	if err != nil {
		act.helper.logValidationErrorMsg(act, err)
		return err
	}

	act.helper.logValidationSuccessMsg(act)
	return nil
}

func (act *BrewAction) Public() bool {
	return true
}

func (act *BrewAction) Run() error {
	act.helper.logRunMsg(act)
	if err := act.helper.runDependencies(act.Deps()); err != nil {
		act.helper.logRunErrorMsg(act, err)
		return err
	}

	for _, pkg := range act.packages {
		err := act.brew.EnsureInstalled(pkg)
		if err != nil {
			act.helper.logRunErrorMsg(act, err)
			return err
		}
	}

	act.helper.logRunSuccessMsg(act)
	return nil
}
