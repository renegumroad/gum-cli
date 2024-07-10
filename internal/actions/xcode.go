package actions

import (
	"slices"

	"github.com/pkg/errors"
	"github.com/renegumroad/gum-cli/internal/cli/cmdexec"
	"github.com/renegumroad/gum-cli/internal/cli/xcode"
	"github.com/renegumroad/gum-cli/internal/systeminfo"
)

type XcodeAction struct {
	cmdGen cmdexec.CmdGenerator
	sys    systeminfo.Client
	xcode  xcode.Client
}

func NewXcodeAction() *XcodeAction {
	return newXcodeActionWithComponents(
		cmdexec.NewCommandGenerator(),
		systeminfo.New(),
		xcode.New(),
	)
}

func newXcodeActionWithComponents(
	gen cmdexec.CmdGenerator,
	sys systeminfo.Client,
	xcode xcode.Client,
) *XcodeAction {
	return &XcodeAction{
		cmdGen: gen,
		sys:    sys,
		xcode:  xcode,
	}
}

func (a *XcodeAction) Name() string {
	return "xcode"
}

func (a *XcodeAction) Identifier() string {
	return "xcode"
}

func (a *XcodeAction) IsPublic() bool {
	return false
}

func (a *XcodeAction) Deps() []Action {
	return []Action{}
}

func (a *XcodeAction) Platforms() []systeminfo.Platform {
	return []systeminfo.Platform{systeminfo.Darwin}
}

func (a *XcodeAction) Validate() error {
	if slices.Contains(a.Platforms(), a.sys.CurrentPlatform()) {
		return nil
	}

	return errors.Errorf("Action %s is not supported on platform %s", a.Name(), a.sys.CurrentPlatform())
}

func (a *XcodeAction) ShouldRun() bool {
	return !a.xcode.IsInstalled()
}

func (a *XcodeAction) Run() error {
	return a.xcode.EnsureInstalled()
}
