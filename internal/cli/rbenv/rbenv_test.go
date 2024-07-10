package rbenv

import (
	"errors"
	"testing"

	"github.com/renegumroad/gum-cli/internal/cli/cmdexec/fakecmdexec"
	"github.com/renegumroad/gum-cli/internal/cli/homebrew"
	"github.com/renegumroad/gum-cli/internal/cli/homebrew/mockhomebrew"
	"github.com/renegumroad/gum-cli/internal/log"
	"github.com/stretchr/testify/suite"
)

type rbenvSuite struct {
	suite.Suite
	mockBrew *mockhomebrew.MockClient
}

func (s *rbenvSuite) SetupSuite() {
	err := log.Initialize(log.LogDisabled)
	s.Require().NoError(err)
}

func (s *rbenvSuite) SetupTest() {
	s.mockBrew = mockhomebrew.NewMockClient(s.T())
}

// TestIsRubyInstalled_Success
func (s *rbenvSuite) TestIsRubyInstalledSuccess() {
	rbenvCmd := fakecmdexec.NewNoOpCommandWithOutputs(&fakecmdexec.NoOpOutputs{
		Stdout: "2.7.2 (set by path)",
	})
	client := newClientWithComponents(fakecmdexec.NewCmdGenerator(rbenvCmd), s.mockBrew)

	result := client.IsRubyInstalled()

	s.Require().True(result)
	s.Require().Equal("rbenv", rbenvCmd.Cmd())
	s.Require().Equal([]string{"version"}, rbenvCmd.Args())
}

func (s *rbenvSuite) TestIsRubyInstalledNotInstalled() {
	rbenvCmd := fakecmdexec.NewNoOpCommandWithOutputs(&fakecmdexec.NoOpOutputs{
		Stdout: "rbenv: version '2.7.2' is not installed (set by path)",
	})
	client := newClientWithComponents(fakecmdexec.NewCmdGenerator(rbenvCmd), s.mockBrew)

	result := client.IsRubyInstalled()

	s.Require().False(result)
	s.Require().Equal("rbenv", rbenvCmd.Cmd())
	s.Require().Equal([]string{"version"}, rbenvCmd.Args())
}

func (s *rbenvSuite) TestIsRubyInstalleErrorRunningCommand() {
	rbenvCmd := fakecmdexec.NewNoOpCommandWithOutputs(&fakecmdexec.NoOpOutputs{
		Err: errors.New("error running command"),
	})
	client := newClientWithComponents(fakecmdexec.NewCmdGenerator(rbenvCmd), s.mockBrew)

	result := client.IsRubyInstalled()

	s.Require().False(result)
	s.Require().Equal("rbenv", rbenvCmd.Cmd())
	s.Require().Equal([]string{"version"}, rbenvCmd.Args())
}

func (s *rbenvSuite) TestEnsureRubyInstalledAlreadyInstalled() {
	rbenvCmdVersion := fakecmdexec.NewNoOpCommandWithOutputs(&fakecmdexec.NoOpOutputs{
		Stdout: "2.7.2 (set by path)",
	})
	client := newClientWithComponents(fakecmdexec.NewCmdGenerator(rbenvCmdVersion), s.mockBrew)

	err := client.EnsureRubyInstalled()

	s.Require().NoError(err)
	s.Require().Equal("rbenv", rbenvCmdVersion.Cmd())
	s.Require().Equal([]string{"version"}, rbenvCmdVersion.Args())
}

func (s *rbenvSuite) TestEnsureRubyInstalledNotInstalled() {
	rbenvCmdVersion := fakecmdexec.NewNoOpCommandWithOutputs(&fakecmdexec.NoOpOutputs{
		Stdout: "2.7.2 not installed (set by path)",
	})
	rbenvInstallCmd := fakecmdexec.NewNoOpCommand()
	s.mockBrew.EXPECT().Upgrade(homebrew.Package{Name: "ruby-build"}).Return(nil)
	client := newClientWithComponents(fakecmdexec.NewCmdGenerator(rbenvCmdVersion, rbenvInstallCmd), s.mockBrew)

	err := client.EnsureRubyInstalled()

	s.Require().NoError(err)
	s.Require().Equal("rbenv", rbenvCmdVersion.Cmd())
	s.Require().Equal([]string{"version"}, rbenvCmdVersion.Args())
	s.Require().Equal("rbenv", rbenvInstallCmd.Cmd())
	s.Require().Equal([]string{"install", "--skip-existing"}, rbenvInstallCmd.Args())
	s.mockBrew.AssertExpectations(s.T())
}

func TestRbenvSuite(t *testing.T) {
	suite.Run(t, new(rbenvSuite))
}
