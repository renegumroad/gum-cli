package xcode

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/renegumroad/gum-cli/internal/cli/cmdexec/fakecmdexec"
	"github.com/renegumroad/gum-cli/internal/log"
	"github.com/stretchr/testify/suite"
)

type xcodeSuite struct {
	suite.Suite
}

func (s *xcodeSuite) SetupSuite() {
	err := log.Initialize(log.LogDisabled)
	s.Require().NoError(err)
}

func (s *xcodeSuite) TestIsInstalled() {
	xcodeCmd := fakecmdexec.NewNoOpCommandWithOutputs(&fakecmdexec.NoOpOutputs{})
	client := newClientWithComponents(fakecmdexec.NewCmdGenerator(xcodeCmd))

	result := client.IsInstalled()

	s.Require().True(result)
	s.Require().Equal("xcode-select", xcodeCmd.Cmd())
	s.Require().Equal([]string{"-p"}, xcodeCmd.Args())
}

func (s *xcodeSuite) TestIsInstalledFalse() {
	xcodeCmd := fakecmdexec.NewNoOpCommandWithOutputs(&fakecmdexec.NoOpOutputs{
		Err: errors.Errorf("xcode-select command not found"),
	})
	client := newClientWithComponents(fakecmdexec.NewCmdGenerator(xcodeCmd))

	result := client.IsInstalled()

	s.Require().False(result)
	s.Require().Equal("xcode-select", xcodeCmd.Cmd())
	s.Require().Equal([]string{"-p"}, xcodeCmd.Args())
}

func (s *xcodeSuite) TestEnsureInstalledAlreadyInstalled() {
	xcodeCheckCmd := fakecmdexec.NewNoOpCommandWithOutputs(&fakecmdexec.NoOpOutputs{
		Err: errors.Errorf("xcode-select command not found"),
	})
	xcodeInstallCmd := fakecmdexec.NewNoOpCommand()

	client := newClientWithComponents(fakecmdexec.NewCmdGenerator(xcodeCheckCmd, xcodeInstallCmd))

	err := client.EnsureInstalled()

	s.Require().NoError(err)
	s.Require().Equal("xcode-select", xcodeCheckCmd.Cmd())
	s.Require().Equal([]string{"-p"}, xcodeCheckCmd.Args())
	s.Require().Equal("xcode-select", xcodeInstallCmd.Cmd())
	s.Require().Equal([]string{"--install"}, xcodeInstallCmd.Args())
}

func (s *xcodeSuite) TestEnsureInstalledNotInstalled() {
	xcodeCheckCmd := fakecmdexec.NewNoOpCommand()
	xcodeInstallCmd := fakecmdexec.NewNoOpCommand()

	client := newClientWithComponents(fakecmdexec.NewCmdGenerator(xcodeCheckCmd, xcodeInstallCmd))

	err := client.EnsureInstalled()

	s.Require().NoError(err)
	s.Require().Equal("xcode-select", xcodeCheckCmd.Cmd())
	s.Require().Equal([]string{"-p"}, xcodeCheckCmd.Args())
	s.Require().Equal("", xcodeInstallCmd.Cmd())
	s.Require().Equal([]string{}, xcodeInstallCmd.Args())
}

func TestXcodeSuite(t *testing.T) {
	suite.Run(t, new(xcodeSuite))
}
