package shellmanager

import (
	"testing"

	"github.com/renehernandez/gum-cli/internal/utils/filesystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type shellmanagerSuite struct {
	suite.Suite
}

type fakeFileSystem struct {
	filesystem.Client
}

func (m *fakeFileSystem) HomeDir() (string, error) {
	return "/home/testuser", nil
}

func (s *shellmanagerSuite) TestGetShellProfilePath() {
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
			actualPath := client.GetShellProfilePath(t.shell)
			assert.Equal(s.T(), t.expectedPath, actualPath)
		})

	}
}

// func TestGetShellProfilePathForSh() {
// 	mockFS := new(MockFileSystem)
// 	mockFS.On("UserHomeDir").Return("/home/testuser", nil)
// 	client := newWithComponents(mockFS)
// 	expectedPath := "/home/testuser/.profile"
// 	actualPath := client.GetShellProfilePath(ShellSh)
// 	assert.Equal(t, expectedPath, actualPath)
// }

// func TestGetShellProfilePathWithUserHomeDirError() {
// 	mockFS := new(MockFileSystem)
// 	mockFS.On("UserHomeDir").Return("", os.ErrNotExist)
// 	client := newWithComponents(mockFS)
// 	actualPath := client.GetShellProfilePath(ShellZsh)
// 	assert.Equal(t, "", actualPath)
// }

func TestShellManagerSuite(t *testing.T) {
	suite.Run(t, new(shellmanagerSuite))
}
