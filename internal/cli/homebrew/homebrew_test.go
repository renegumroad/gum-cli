package homebrew

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/renegumroad/gum-cli/internal/cli/cmdexec"
	"github.com/renegumroad/gum-cli/internal/cli/cmdexec/fakecmdexec"
	"github.com/renegumroad/gum-cli/internal/filesystem"
	"github.com/renegumroad/gum-cli/internal/filesystem/mockfilesystem"
	"github.com/renegumroad/gum-cli/internal/log"
	"github.com/stretchr/testify/suite"
)

type brewSuite struct {
	suite.Suite
	mockFs         *mockfilesystem.MockClient
	client         Client
	origBrewPrefix string
	testBrewPrefix string
	pkgPath        string
	caskPath       string
}

func (s *brewSuite) SetupSuite() {
	err := log.Initialize(log.LogDisabled)
	s.Require().NoError(err)
}

func (s *brewSuite) SetupTest() {
	s.mockFs = mockfilesystem.NewMockClient(s.T())
	s.client = newClientWithComponents(s.mockFs, cmdexec.NewCommandGenerator())

	s.origBrewPrefix = os.Getenv("HOMEBREW_PREFIX")

	fs := filesystem.New()

	path, err := fs.MkdirTemp()
	s.Require().NoError(err)

	s.testBrewPrefix = path
	os.Setenv("HOMEBREW_PREFIX", path)
	s.pkgPath = filepath.Join(path, "opt")
	s.caskPath = filepath.Join(path, "Caskroom")
}

func (s *brewSuite) TearDownTest() {
	os.Setenv("HOMEBREW_PREFIX", s.origBrewPrefix)
}

func (s *brewSuite) TestIsInstalledTrue() {
	pkg := Package{Name: "testpkg"}

	s.mockFs.EXPECT().Exists(filepath.Join(s.pkgPath, pkg.Name)).Return(true)

	s.Require().True(s.client.IsInstalled(pkg))
}

func (s *brewSuite) TestIsInstalledFalse() {
	pkg := Package{Name: "testpkg"}
	s.mockFs.EXPECT().Exists(filepath.Join(s.pkgPath, pkg.Name)).Return(false)

	installed := s.client.IsInstalled(pkg)
	s.Require().False(installed)
}

func (s *brewSuite) TestIsInstalledWhenNameIsEmpty() {
	pkg := Package{Name: ""}
	s.Require().False(s.client.IsInstalled(pkg))
}

func (s *brewSuite) TestIsInstalledCaskTrue() {
	pkg := Package{Name: "testpkg", Cask: true}

	s.mockFs.EXPECT().Exists(filepath.Join(s.caskPath, pkg.Name)).Return(true)

	s.Require().True(s.client.IsInstalled(pkg))
}

func (s *brewSuite) TestIsInstalledCaskFalse() {
	pkg := Package{Name: "testpkg", Cask: true}
	s.mockFs.EXPECT().Exists(filepath.Join(s.caskPath, pkg.Name)).Return(false)

	installed := s.client.IsInstalled(pkg)
	s.Require().False(installed)
}

func (s *brewSuite) TestInstall() {
	pkg := Package{Name: "testpkg"}
	noOpCmd := fakecmdexec.NewNoOpCommand()
	s.client = newClientWithComponents(s.mockFs, fakecmdexec.NewCmdGenerator(noOpCmd))

	err := s.client.Install(pkg)
	s.Require().NoError(err)

	s.Require().Equal("brew", noOpCmd.Cmd())
	s.Require().Equal([]string{"install", "testpkg"}, noOpCmd.Args())
}

func (s *brewSuite) TestInstallCask() {
	pkg := Package{Name: "testpkg", Cask: true}
	noOpCmd := fakecmdexec.NewNoOpCommand()
	s.client = newClientWithComponents(s.mockFs, fakecmdexec.NewCmdGenerator(noOpCmd))

	err := s.client.Install(pkg)
	s.Require().NoError(err)

	s.Require().Equal("brew", noOpCmd.Cmd())
	s.Require().Equal([]string{"install", "--cask", "testpkg"}, noOpCmd.Args())
}

func (s *brewSuite) TestLink() {
	pkg := Package{Name: "testpkg", Link: true}
	noOpCmd := fakecmdexec.NewNoOpCommand()
	s.client = newClientWithComponents(s.mockFs, fakecmdexec.NewCmdGenerator(noOpCmd))

	err := s.client.Link(pkg)
	s.Require().NoError(err)

	s.Require().Equal("brew", noOpCmd.Cmd())
	s.Require().Equal([]string{"link", "--force", "--overwrite", "testpkg"}, noOpCmd.Args())
}

func (s *brewSuite) TestLinkWhenNoName() {
	pkg := Package{Name: "", Link: true}
	noOpCmd := fakecmdexec.NewNoOpCommand()
	s.client = newClientWithComponents(s.mockFs, fakecmdexec.NewCmdGenerator(noOpCmd))

	err := s.client.Link(pkg)

	s.Require().Error(err)
	s.Require().Equal("", noOpCmd.Cmd())
	s.Require().Equal([]string{}, noOpCmd.Args())
}

func (s *brewSuite) TestLinkCask() {
	pkg := Package{Name: "testpkg", Cask: true, Link: true}
	noOpCmd := fakecmdexec.NewNoOpCommand()
	s.client = newClientWithComponents(s.mockFs, fakecmdexec.NewCmdGenerator(noOpCmd))

	err := s.client.Link(pkg)

	s.Require().Error(err)
	s.Require().Equal("Cannot link cask package testpkg", err.Error())
	s.Require().Equal("", noOpCmd.Cmd())
	s.Require().Equal([]string{}, noOpCmd.Args())
}

func (s *brewSuite) TestLinkNotLink() {
	pkg := Package{Name: "testpkg"}
	noOpCmd := fakecmdexec.NewNoOpCommand()
	s.client = newClientWithComponents(s.mockFs, fakecmdexec.NewCmdGenerator(noOpCmd))

	err := s.client.Link(pkg)
	s.Require().NoError(err)
	s.Require().Equal("", noOpCmd.Cmd())
	s.Require().Equal([]string{}, noOpCmd.Args())
}

func (s *brewSuite) TestEnsureInstalledAlreadyInstalled() {
	pkg := Package{Name: "testpkg"}
	s.mockFs.EXPECT().Exists(filepath.Join(s.pkgPath, pkg.Name)).Return(true)

	noOpCmd := fakecmdexec.NewNoOpCommand()
	s.client = newClientWithComponents(s.mockFs, fakecmdexec.NewCmdGenerator(noOpCmd))

	err := s.client.EnsureInstalled(pkg)
	s.Require().NoError(err)
	s.Require().Equal("", noOpCmd.Cmd())
	s.Require().Equal([]string{}, noOpCmd.Args())
}

func (s *brewSuite) TestEnsureInstalledNotInstalled() {
	pkg := Package{Name: "testpkg"}
	s.mockFs.EXPECT().Exists(filepath.Join(s.pkgPath, pkg.Name)).Return(false)

	noOpCmd := fakecmdexec.NewNoOpCommand()
	s.client = newClientWithComponents(s.mockFs, fakecmdexec.NewCmdGenerator(noOpCmd))
	err := s.client.EnsureInstalled(pkg)
	s.Require().NoError(err)
	s.Require().Equal("brew", noOpCmd.Cmd())
	s.Require().Equal([]string{"install", "testpkg"}, noOpCmd.Args())
}

func TestBrewSuite(t *testing.T) {
	suite.Run(t, &brewSuite{})
}
