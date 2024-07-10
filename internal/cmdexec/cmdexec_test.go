package cmdexec

import (
	"strings"
	"testing"

	"github.com/renegumroad/gum-cli/internal/log"
	"github.com/stretchr/testify/suite"
)

type cmdExecSuite struct {
	suite.Suite
}

func (suite *cmdExecSuite) SetupSuite() {
	err := log.Initialize(log.LogDisabled)
	suite.Require().NoError(err)
}

func (suite *cmdExecSuite) TestRunSuccess() {
	cmd := New("echo", "hello")
	err := cmd.Run()
	suite.NoError(err, "Expected no error")
	suite.Equal("hello", strings.TrimSpace(cmd.Stdout), "Stdout should match expected output")
	suite.Empty(cmd.Stderr, "Stderr should be empty")
}

func (suite *cmdExecSuite) TestRunWithEnvVar() {
	cmd := NewWithEnv("printenv", []string{"TEST_VAR"}, []string{"TEST_VAR=value"})
	err := cmd.Run()
	suite.NoError(err, "Expected no error")
	suite.Equal("value", strings.TrimSpace(cmd.Stdout), "Stdout should contain environment variable value")
	suite.Empty(cmd.Stderr, "Stderr should be empty")
}

func (suite *cmdExecSuite) TestRunFailure() {
	cmd := New("false")
	err := cmd.Run()
	suite.Error(err, "Expected an error")
	suite.Empty(cmd.Stdout, "Stdout should be empty")
}

func (suite *cmdExecSuite) TestRunCommandNotFound() {
	cmd := New("nonexistentcommand")
	err := cmd.Run()
	suite.Error(err, "Expected an error")
}

func TestCmdExecSuite(t *testing.T) {
	suite.Run(t, new(cmdExecSuite))
}
