package shellmanager

import (
	"os"
	"testing"

	"github.com/renegumroad/gum-cli/internal/filesystem"
	"github.com/stretchr/testify/suite"
)

type shellManagerSuite struct {
	suite.Suite
}

type fakeFileSystem struct {
	filesystem.Client
}

func (m *fakeFileSystem) HomeDir() (string, error) {
	return "/home/testuser", nil
}

type invalidHomeDirFileSystem struct {
	filesystem.Client
}

func (m *invalidHomeDirFileSystem) HomeDir() (string, error) {
	return "", os.ErrNotExist
}

func (s *shellManagerSuite) TestGetShellProfilePath() {
	tests := []struct {
		name         string
		shell        ShellType
		expectedPath string
	}{
		{
			name:         "Zsh",
			shell:        ShellZsh,
			expectedPath: "/home/testuser/.zprofile",
		},
		{
			name:         "Bash",
			shell:        ShellBash,
			expectedPath: "/home/testuser/.bash_profile",
		},
		{
			name:         "Sh",
			shell:        ShellSh,
			expectedPath: "/home/testuser/.profile",
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			client := newWithComponents(&fakeFileSystem{})
			actualPath, err := client.GetShellProfilePath(t.shell)
			s.Require().NoError(err)
			s.Require().Equal(t.expectedPath, actualPath)
		})

	}
}

func (s *shellManagerSuite) TestGetShellProfilePathWithUserHomeDirError() {
	client := newWithComponents(&invalidHomeDirFileSystem{})
	actualPath, err := client.GetShellProfilePath(ShellZsh)
	s.Require().Error(err)
	s.Require().Empty(actualPath)
}

func TestShellManagerSuite(t *testing.T) {
	suite.Run(t, new(shellManagerSuite))
}
