package projectconfig

import (
	"testing"

	"github.com/renegumroad/gum-cli/internal/github/mockgithub"
	"github.com/renegumroad/gum-cli/internal/log"
	"github.com/stretchr/testify/suite"
)

type projectConfigSuite struct {
	suite.Suite
	mockGithub *mockgithub.MockClient
}

func (suite *projectConfigSuite) SetupSuite() {
	err := log.Initialize(log.LogDisabled)
	suite.Require().NoError(err)
}

func (s *projectConfigSuite) SetupTest() {
	s.mockGithub = mockgithub.NewMockClient(s.T())
}

func (s *projectConfigSuite) TestExtractOwnerAndRepoFromURL() {
	validURL := "https://github.com/renegumroad/gum-cli"
	expectedOwner := "renegumroad"
	expectedRepo := "gum-cli"
	owner, repo, err := extractOwnerAndRepoFromURL(validURL)
	s.Require().NoError(err)
	s.Require().Equal(expectedOwner, owner)
	s.Require().Equal(expectedRepo, repo)
}

func (s *projectConfigSuite) TestExtractOwnerAndRepoFromURLInsufficientPathSegments() {
	insufficientURL := "https://github.com/renegumroad"
	_, _, err := extractOwnerAndRepoFromURL(insufficientURL)
	s.Require().Error(err)
}

func (s *projectConfigSuite) TestNewWithComponentsValidArgsWithNameAndURL() {
	args := &ProjectConfigArgs{Name: "gum-cli", URL: "https://github.com/renegumroad/gum-cli", Owner: "renegumroad"}
	project, err := newWithComponents(args, s.mockGithub)
	s.Require().NoError(err)
	s.Require().Equal("gum-cli", project.Name)
	s.Require().Equal("https://github.com/renegumroad/gum-cli", project.URL)
	s.Require().Equal("renegumroad", project.Owner)
}

func (s *projectConfigSuite) TestNewWithComponentsValidArgsWithURLOnly() {
	args := &ProjectConfigArgs{URL: "https://github.com/renegumroad/gum-cli"}
	project, err := newWithComponents(args, s.mockGithub)
	s.Require().NoError(err)
	s.Require().Equal("gum-cli", project.Name)
}

func (s *projectConfigSuite) TestNewWithComponentsValidArgsWithNameOnly() {
	args := &ProjectConfigArgs{Name: "gum-cli", Owner: "renegumroad"}
	project, err := newWithComponents(args, s.mockGithub)
	s.Require().NoError(err)
	s.Require().Equal("https://github.com/renegumroad/gum-cli", project.URL)
}

func (s *projectConfigSuite) TestNewWithComponentsInvalidArgsNoNameOrURL() {
	args := &ProjectConfigArgs{}
	_, err := newWithComponents(args, s.mockGithub)
	s.Require().Error(err)
}

func (s *projectConfigSuite) TestNewWithComponentsInvalidURL() {
	args := &ProjectConfigArgs{URL: "https://invalid.com/renegumroad/gum-cli"}
	_, err := newWithComponents(args, s.mockGithub)
	s.Require().Error(err)
}

func (s *projectConfigSuite) TestNewWithComponentsMismatchNameAndURL() {
	args := &ProjectConfigArgs{Name: "different-name", URL: "https://github.com/renegumroad/gum-cli"}
	_, err := newWithComponents(args, s.mockGithub)
	s.Require().Error(err)
}

func (s *projectConfigSuite) TestNewWithComponentsMismatchOwnerAndURL() {
	args := &ProjectConfigArgs{Owner: "different-owner", URL: "https://github.com/renegumroad/gum-cli"}
	_, err := newWithComponents(args, s.mockGithub)
	s.Require().Error(err)
}

func TestProjectConfigSuite(t *testing.T) {
	suite.Run(t, new(projectConfigSuite))
}
