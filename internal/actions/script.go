package actions

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/renegumroad/gum-cli/internal/cli/cmdexec"
	"github.com/renegumroad/gum-cli/internal/log"
	"github.com/renegumroad/gum-cli/internal/systeminfo"
)

type ScriptActionArgs struct {
	Title   string
	Test    string
	Command string
}

type ScriptAction struct {
	args   *ScriptActionArgs
	cmdGen cmdexec.CmdGenerator
}

func NewScriptAction(args *ScriptActionArgs) *ScriptAction {
	return newScriptActionWithComponents(args, cmdexec.NewCommandGenerator())
}

func newScriptActionWithComponents(args *ScriptActionArgs, gen cmdexec.CmdGenerator) *ScriptAction {
	return &ScriptAction{
		args:   args,
		cmdGen: gen,
	}
}

func (a *ScriptAction) Name() string {
	return "script"
}

func (a *ScriptAction) Identifier() string {
	return a.args.Title
}

func (a *ScriptAction) IsPublic() bool {
	return true
}

func (a *ScriptAction) Platforms() []systeminfo.Platform {
	return []systeminfo.Platform{systeminfo.Darwin, systeminfo.Linux}
}

func (a *ScriptAction) Deps() []Action {
	return []Action{}
}

func (a *ScriptAction) Validate() error {
	errMsg := fmt.Sprintf("Failed %s action validation", a.Name())
	errFound := false

	if a.args.Title == "" {
		errMsg = fmt.Sprintf("%s: title is empty", errMsg)
		errFound = true
	}

	if a.args.Command == "" {
		errMsg = fmt.Sprintf("%s: command is empty", errMsg)
		errFound = true
	}

	if errFound {
		return errors.Errorf(errMsg)
	}
	return nil
}

func (a *ScriptAction) ShouldRun() bool {
	if a.args.Test != "" {
		if err := a.runCmd(a.args.Test); err == nil {
			log.Debugf("Script test %s passed", a.args.Test)
			return false
		} else {
			log.Debugf("Script test %s failed", a.args.Test)
		}
	}
	return true
}

func (a *ScriptAction) Run() error {

	return a.runCmd(a.args.Command)
}

func (a *ScriptAction) runCmd(cmd string) error {
	err := a.cmdGen("bash", "-c", cmd).Run()
	if err != nil {
		return errors.Errorf("Failed to run command: %s", err)
	}

	return nil
}
