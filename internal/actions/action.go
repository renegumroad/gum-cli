package actions

import (
	"slices"

	"github.com/pkg/errors"
	"github.com/renegumroad/gum-cli/internal/log"
	"github.com/renegumroad/gum-cli/internal/systeminfo"
)

var (
	namedActions = map[string]Action{
		"golang":      NewGolangAction(),
		"ruby":        NewRubyAction(),
		"xcode":       NewXcodeAction(),
		"brew_ensure": NewBrewEnsureAction(),
	}
)

type Platform string

type Action interface {
	Name() string
	Identifier() string
	Deps() []Action
	Validate() error
	ShouldRun() bool
	Run() error
	IsPublic() bool
	Platforms() []systeminfo.Platform
}

func SupportedByConfig(name string) bool {
	return namedActions[name] != nil && namedActions[name].IsPublic()
}

func Get(name string) Action {
	return namedActions[name]
}

func SupportedByCurrentPlatform(action Action) bool {
	sys := systeminfo.New()

	return slices.Contains(action.Platforms(), sys.CurrentPlatform())
}

type ActionHandler struct {
	Actions []Action
}

func NewActionHandler(actions []Action) *ActionHandler {
	return &ActionHandler{
		Actions: buildActionList(actions...),
	}
}

func (h *ActionHandler) Validate() error {
	errMsg := "Failed action(s) validation(s):"
	errFound := false

	for _, action := range h.Actions {
		log.Debugf("Validating action %s", action.Name())

		if err := action.Validate(); err != nil {
			log.Debugf("Action %s failed validation", action.Name())
			errMsg = errMsg + "\n" + err.Error()
			errFound = true
		} else {
			log.Debugf("Action %s validated", action.Name())
		}
	}

	if errFound {
		return errors.Errorf(errMsg)
	}
	return nil
}

func (h *ActionHandler) Run() error {
	for _, action := range h.Actions {
		if !action.ShouldRun() {
			log.Infof("Skipping action %s", action.Name())
			continue
		}
		log.Infof("Running action %s", action.Name())
		if err := action.Run(); err != nil {
			log.Errorf("Action %s run failed: %s", action.Name(), err)
			return err
		} else {
			log.Infof("Action %s ran successfully", action.Name())
		}
	}

	return nil
}

func buildActionList(actions ...Action) []Action {
	sortedActions := []Action{}

	for _, action := range actions {
		if !SupportedByCurrentPlatform(action) {
			log.Debugf("Skipping %s action. Not supported by current platform", action.Name())
			continue
		}

		depsActions := buildActionList(action.Deps()...)

		for _, depAction := range depsActions {
			if slices.ContainsFunc(sortedActions, containsAction(depAction)) {
				continue
			}

			sortedActions = append(sortedActions, depAction)
		}

		if slices.ContainsFunc(sortedActions, containsAction(action)) {
			continue
		}
		sortedActions = append(sortedActions, action)
	}

	return sortedActions
}

func containsAction(a Action) func(b Action) bool {
	return func(b Action) bool {
		return a.Identifier() == b.Identifier()
	}
}

func depsShouldRun(actions []Action) bool {
	for _, action := range actions {
		if action.ShouldRun() {
			return true
		}
	}

	return false
}
