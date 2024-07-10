package actions

import (
	"github.com/pkg/errors"
	"github.com/renegumroad/gum-cli/internal/log"
)

var (
	namedActions = map[string]Action{
		"golang": &GolangAction{},
		"ruby":   &RubyAction{},
	}
)

type Action interface {
	Name() string
	Deps() []Action
	Validate() error
	Run() error
	Public() bool
}

func SupportedByConfig(name string) bool {
	return namedActions[name] != nil && namedActions[name].Public()
}

func Get(name string) Action {
	return namedActions[name]
}

func Validate(action Action) error {
	helper := &actionHelper{}

	log.Debugf("Validating action %s", action.Name())
	depsError := helper.validateDependencies(action.Deps())

	err := action.Validate()

	if depsError != nil {
		if err != nil {
			err = errors.Errorf("%s: deps: %s", err, depsError)
		} else {
			err = depsError
		}
	}

	if err != nil {
		log.Debugf("Action %s failed validation", action.Name())
	} else {
		log.Debugf("Action %s validated", action.Name())
	}

	return err
}

func Run(action Action) error {
	helper := &actionHelper{}

	log.Infof("Running action %s", action.Name())
	if err := helper.runDependencies(action.Deps()); err != nil {
		return err
	}

	if err := action.Run(); err != nil {
		log.Errorf("Action %s run failed: %s", action.Name(), err)
		return err
	}

	log.Infof("Action %s ran successfully", action.Name())

	return nil
}

type actionHelper struct {
}

func (h *actionHelper) validateDependencies(deps []Action) error {
	var err error
	for _, dep := range deps {
		log.Debugf("Validating dependent action %s", dep.Name())
		depError := dep.Validate()
		if depError != nil {
			err = errors.Errorf("%s\n- %s: %s", dep.Name(), depError, err)
		} else {
			log.Debugf("Dependent action %s validated", dep.Name())
		}
	}

	if err != nil {
		err = errors.Errorf("Failed dependent actions validation:\n%s", err)
	}

	return err
}

func (h *actionHelper) runDependencies(deps []Action) error {
	for _, dep := range deps {
		log.Debugf("Running dependent action %s", dep.Name())
		err := dep.Run()
		if err != nil {
			return errors.Errorf("Failed running dependent action %s: %s", dep.Name(), err)
		}

		log.Debugf("Dependent action %s ran successfully", dep.Name())
	}

	return nil
}
