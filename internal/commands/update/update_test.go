package update

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/renegumroad/gum-cli/internal/log"
	"github.com/renegumroad/gum-cli/internal/projectconfig/mockprojectconfig"
	"github.com/renegumroad/gum-cli/internal/version"
	"github.com/stretchr/testify/suite"
)

type updateImplTestSuite struct {
	suite.Suite
	updateImpl        *UpdateImpl
	mockProjectConfig *mockprojectconfig.MockProjectConfig
	codeVersion       string
}

func (suite *updateImplTestSuite) SetupSuite() {
	err := log.Initialize(log.LogDisabled)
	suite.Require().NoError(err)
}

func (s *updateImplTestSuite) SetupTest() {
	s.codeVersion = version.VERSION
	s.mockProjectConfig = mockprojectconfig.NewMockProjectConfig(s.T())
	s.updateImpl = newWithComponents(s.mockProjectConfig)
}

func (s *updateImplTestSuite) TearDownTest() {
	version.VERSION = s.codeVersion
}

func (s *updateImplTestSuite) TestNoUpdateNeeded() {
	version.VERSION = "0.0.6"
	s.mockProjectConfig.On("LatestVersion").Return(version.VERSION, nil)
	needsUpdate, err := s.updateImpl.checkforUpdate()
	s.False(needsUpdate)
	s.NoError(err)
}

func (s *updateImplTestSuite) TestUpdateNeeded() {
	version.VERSION = "0.0.6"
	s.mockProjectConfig.On("LatestVersion").Return("0.1.0", nil) // Assuming this is greater than current version.VERSION
	needsUpdate, err := s.updateImpl.checkforUpdate()
	s.True(needsUpdate)
	s.NoError(err)
}

func (s *updateImplTestSuite) TestInvalidCurrentVersionSemver() {
	version.VERSION = "invalid.semver"
	_, err := s.updateImpl.checkforUpdate()
	s.Error(err)
}

func (s *updateImplTestSuite) TestInvalidLatestVersionSemver() {
	version.VERSION = "0.0.6"
	s.mockProjectConfig.On("LatestVersion").Return("invalid.semver", nil)
	_, err := s.updateImpl.checkforUpdate()
	s.Error(err)
}

func (s *updateImplTestSuite) TestLatestVersionFetchError() {
	version.VERSION = "0.0.6"
	s.mockProjectConfig.On("LatestVersion").Return("", errors.Errorf("error fetching latest version"))
	_, err := s.updateImpl.checkforUpdate()
	s.Error(err)
}

func TestUpdateImplTestSuite(t *testing.T) {
	suite.Run(t, new(updateImplTestSuite))
}
