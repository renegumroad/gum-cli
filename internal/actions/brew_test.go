package actions

import (
	"testing"

	"github.com/renegumroad/gum-cli/internal/cli/homebrew"
	"github.com/renegumroad/gum-cli/internal/cli/homebrew/mockhomebrew"
	"github.com/renegumroad/gum-cli/internal/log"
	"github.com/stretchr/testify/suite"
)

type brewActionSuite struct {
	suite.Suite
	mockBrew *mockhomebrew.MockClient
}

func (suite *brewActionSuite) SetupSuite() {
	err := log.Initialize(log.LogDisabled)
	suite.Require().NoError(err)
}

func (s *brewActionSuite) SetupTest() {
	s.mockBrew = mockhomebrew.NewMockClient(s.T())
}

func (s *brewActionSuite) TestValidateNoPackagesError() {
	act := NewBrewAction([]homebrew.Package{})

	err := act.Validate()
	s.Require().ErrorContains(err, "no packages specified")
}

func (s *brewActionSuite) TestValidateSomePackagesDoNotHaveName() {
	act := NewBrewAction([]homebrew.Package{{Name: ""}, {Name: "package1"}})

	err := act.Validate()
	s.Require().ErrorContains(err, "package(s) missing name")
}

func (s *brewActionSuite) TestRun() {
	pkgs := []homebrew.Package{{Name: "package1"}, {Name: "package2"}}
	for _, pkg := range pkgs {
		s.mockBrew.EXPECT().EnsureInstalled(pkg).Return(nil)
	}
	act := newBrewActionWithClient(pkgs, s.mockBrew)

	err := act.Run()
	s.Require().NoError(err)
	s.mockBrew.AssertNumberOfCalls(s.T(), "EnsureInstalled", 2)
}

func TestBrewActionSuite(t *testing.T) {
	suite.Run(t, new(brewActionSuite))
}
