package actions

import (
	"github.com/pkg/errors"
	"github.com/renegumroad/gum-cli/internal/log"
)

var (
	namedActions = map[string]Action{}
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

type actionHelper struct {
}

func (h *actionHelper) validateDependencies(deps []Action) error {
	var err error
	for _, dep := range deps {
		depError := dep.Validate()
		if depError != nil {
			err = errors.Errorf("%s\n- %s: %s", dep.Name(), depError, err)
		}
	}

	if err != nil {
		err = errors.Errorf("Failed dependent actions validation:\n%s", err)
	}

	return err
}

func (h *actionHelper) runDependencies(deps []Action) error {
	for _, dep := range deps {
		err := dep.Run()
		if err != nil {
			return errors.Errorf("Failed running dependent action %s: %s", dep.Name(), err)
		}
	}

	return nil
}

func (h *actionHelper) logValidationMsg(action Action) {
	log.Debugf("Validating action %s", action.Name())
}

func (h *actionHelper) logValidationSuccessMsg(action Action) {
	log.Debugf("Action %s validated successfully", action.Name())
}

func (h *actionHelper) logValidationErrorMsg(action Action, err error) {
	log.Errorf("Action %s validation failed: %s", action.Name(), err)
}

func (h *actionHelper) logRunMsg(action Action) {
	log.Infof("Running action %s", action.Name())
}

func (h *actionHelper) logRunSuccessMsg(action Action) {
	log.Infof("Action %s ran successfully", action.Name())
}

func (h *actionHelper) logRunErrorMsg(action Action, err error) {
	log.Errorf("Action %s run failed: %s", action.Name(), err)
}
