package bundler

import (
	"testing"

	"github.com/renegumroad/gum-cli/internal/filesystem/mockfilesystem"
	"github.com/renegumroad/gum-cli/internal/log"
	"github.com/stretchr/testify/suite"
)

type bundlerSuite struct {
	suite.Suite
	mockFs *mockfilesystem.MockClient
}

func (s *bundlerSuite) SetupSuite() {
	err := log.Initialize(log.LogDisabled)
	s.Require().NoError(err)
}

func (s *bundlerSuite) SetupTest() {
	s.mockFs = mockfilesystem.NewMockClient(s.T())
}

func (s *bundlerSuite) TestGetVersionFromVersionFileExistsWithValidVersion() {
	s.mockFs.On("CurrentDir").Return("/test/dir", nil)
	s.mockFs.On("Exists", "/test/dir/.bundler-version").Return(true)
	s.mockFs.On("ReadString", "/test/dir/.bundler-version").Return("2.1.4", nil)

	client := newClientWithComponents(s.mockFs, nil)

	version, err := client.getVersionFromVersionFile()

	s.Require().NoError(err)
	s.Equal("2.1.4", version)
	s.mockFs.AssertExpectations(s.T())
}

func (s *bundlerSuite) TestGetVersionFromVersionFileNotExists() {
	s.mockFs.On("CurrentDir").Return("/test/dir", nil)
	s.mockFs.On("Exists", "/test/dir/.bundler-version").Return(false)

	client := newClientWithComponents(s.mockFs, nil)

	version, err := client.getVersionFromVersionFile()

	s.Require().NoError(err)
	s.Equal("", version)
	s.mockFs.AssertExpectations(s.T())
}

func (s *bundlerSuite) TestGetVersionFromVersionFileExistsWithInvalidVersion() {
	s.mockFs.On("CurrentDir").Return("/test/dir", nil)
	s.mockFs.On("Exists", "/test/dir/.bundler-version").Return(true)
	s.mockFs.On("ReadString", "/test/dir/.bundler-version").Return("invalid_version", nil)

	client := newClientWithComponents(s.mockFs, nil)

	version, err := client.getVersionFromVersionFile()

	s.Require().Error(err)
	s.Equal("", version)
	s.mockFs.AssertExpectations(s.T())
}

func (s *bundlerSuite) TestGetVersionFromVersionFileExistsWithEmptyVersion() {
	s.mockFs.On("CurrentDir").Return("/test/dir", nil)
	s.mockFs.On("Exists", "/test/dir/.bundler-version").Return(true)
	s.mockFs.On("ReadString", "/test/dir/.bundler-version").Return("", nil)

	client := newClientWithComponents(s.mockFs, nil)

	version, err := client.getVersionFromVersionFile()

	s.Require().Error(err)
	s.Equal("", version)
	s.mockFs.AssertExpectations(s.T())
}

func (s *bundlerSuite) TestGetVersionFromGemfileLockExistsWithValidVersion() {
	s.mockFs.On("CurrentDir").Return("/test/dir", nil)
	s.mockFs.On("Exists", "/test/dir/Gemfile.lock").Return(true)
	s.mockFs.On("ReadString", "/test/dir/Gemfile.lock").Return("BUNDLED WITH\n   2.1.4", nil)

	client := newClientWithComponents(s.mockFs, nil)

	version, err := client.getVersionFromGemfileLock()

	s.Require().NoError(err)
	s.Equal("2.1.4", version)
	s.mockFs.AssertExpectations(s.T())
}

func (s *bundlerSuite) TestGetVersionFromGemfileLockNotExists() {
	s.mockFs.On("CurrentDir").Return("/test/dir", nil)
	s.mockFs.On("Exists", "/test/dir/Gemfile.lock").Return(false)

	client := newClientWithComponents(s.mockFs, nil)

	version, err := client.getVersionFromGemfileLock()

	s.Require().NoError(err)
	s.Equal("", version)
	s.mockFs.AssertExpectations(s.T())
}

func (s *bundlerSuite) TestGetVersionFromGemfileLockExistsWithInvalidVersion() {
	s.mockFs.On("CurrentDir").Return("/test/dir", nil)
	s.mockFs.On("Exists", "/test/dir/Gemfile.lock").Return(true)
	s.mockFs.On("ReadString", "/test/dir/Gemfile.lock").Return("BUNDLED WITH\n   invalid_version", nil)

	client := newClientWithComponents(s.mockFs, nil)

	version, err := client.getVersionFromGemfileLock()

	s.Require().Error(err)
	s.Equal("", version)
	s.mockFs.AssertExpectations(s.T())
}

func (s *bundlerSuite) TestGetVersionFromGemfileLockExistsWithEmptyVersion() {
	s.mockFs.On("CurrentDir").Return("/test/dir", nil)
	s.mockFs.On("Exists", "/test/dir/Gemfile.lock").Return(true)
	s.mockFs.On("ReadString", "/test/dir/Gemfile.lock").Return("", nil)

	client := newClientWithComponents(s.mockFs, nil)

	version, err := client.getVersionFromGemfileLock()

	s.Require().Error(err)
	s.Equal("", version)
	s.mockFs.AssertExpectations(s.T())
}

func TestBundlerSuite(t *testing.T) {
	suite.Run(t, new(bundlerSuite))
}
