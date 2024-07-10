package filesystem

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
	"testing"

	"github.com/stretchr/testify/suite"
)

type filesystemSuite struct {
	suite.Suite
}

func (s *filesystemSuite) TestCurrentDir() {
	c := NewClient()
	path, err := c.CurrentDir()

	s.Require().NoError(err)
	s.Require().NotEmpty(path)
}

func (s *filesystemSuite) TestExists() {
	c := NewClient()
	path, _ := c.CurrentDir()

	path = filepath.Join(path, "filesystem.go")

	s.Require().FileExists(path)
}

func (s *filesystemSuite) TestCreateTempDir() {
	c := NewClient()
	tempDir, err := c.CreateTempDir()

	s.Require().NoError(err, "Expected no error on creating a temp dir")
	s.Require().DirExists(tempDir)

	// Clean up after the test
	defer os.RemoveAll(tempDir)

	// Check that the directory name matches the expected pattern
	s.Require().Regexp(`gum_.*`, filepath.Base(tempDir))
}

func (s *filesystemSuite) TestCopyFile() {
	c := NewClient()
	path, _ := c.CurrentDir()

	dir, err := c.CreateTempDir()
	s.Require().NoError(err)
	s.Require().DirExists(dir)
	defer os.RemoveAll(dir)

	source := filepath.Join(path, "filesystem.go")
	destination := filepath.Join(dir, "filesystem.go")

	s.Require().NoError(c.CopyFile(source, destination))
	s.Require().FileExists(destination)

	equal, err := c.EqualFiles(source, destination)
	s.Require().NoError(err)
	s.Require().True(equal)
}

func (s *filesystemSuite) TestIsSymlink() {
	c := NewClient()
	tempDir, err := c.CreateTempDir()
	s.Require().NoError(err, "Expected no error on creating a temp dir")
	defer os.RemoveAll(tempDir)

	// Create a temp file
	path := filepath.Join(tempDir, "tempfile")
	err = c.WriteString(path, "test")
	s.Require().NoError(err, "Expected no error on writing a temp file")

	// Create a symlink to the temp file
	symlinkPath := filepath.Join(tempDir, "symlink")
	err = os.Symlink(path, symlinkPath)
	s.Require().NoError(err, "Expected no error on creating a symlink")

	s.Require().True(c.IsSymlink(symlinkPath))
	s.Require().False(c.IsSymlink(path))
}

func (s *filesystemSuite) TestWriteString() {
	c := NewClient()
	tmpDir, err := c.CreateTempDir()
	s.Require().NoError(err)
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "testfile.txt")
	content := "Hello, World!"

	// Test writing new content to a new file
	err = c.WriteString(filePath, content)
	s.Require().NoError(err, "Expected no error on writing to a new file")
	s.Require().FileExists(filePath)

	readContent, err := os.ReadFile(filePath)
	s.Require().NoError(err, "Expected no error on reading the file")
	s.Require().Equal(content, string(readContent), "Content written and read should match")

	// Test overwriting existing file
	newContent := "Goodbye, World!"
	err = c.WriteString(filePath, newContent)
	s.Require().NoError(err, "Expected no error on overwriting an existing file")

	readContent, err = os.ReadFile(filePath)
	s.Require().NoError(err, "Expected no error on reading the overwritten file")
	s.Require().Equal(newContent, string(readContent), "Overwritten content should match the new content")
}

func (s *filesystemSuite) TestGetOwner() {
	c := NewClient()
	tmpDir, err := c.CreateTempDir()
	s.Require().NoError(err)
	defer os.RemoveAll(tmpDir)

	// Create a temporary file to test ownership
	tmpFile, err := os.CreateTemp(tmpDir, "testfile-*")
	s.Require().NoError(err)
	defer tmpFile.Close()

	// Test getting owner of the file
	userInfo, err := c.GetOwner(tmpFile.Name())
	s.Require().NoError(err, "Expected no error on getting file owner")

	// Get current user as expected owner
	currentUser, err := user.Current()
	s.Require().NoError(err, "Expected no error on getting current user")

	// Check if the owner matches the current user
	s.Require().Equal(currentUser.Username, userInfo.Name, "The owner username should match the current user")
	s.Require().Equal(currentUser.Uid, fmt.Sprint(userInfo.Id), "The owner user ID should match the current user ID")
}

func (s *filesystemSuite) TestChownRecursively() {
	c := NewClient()
	tmpDir, err := c.CreateTempDir()
	s.Require().NoError(err)
	defer os.RemoveAll(tmpDir)

	// Create a subdirectory and a file inside tmpDir to test recursive chown
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	s.Require().NoError(err)

	tmpFile, err := os.CreateTemp(subDir, "testfile-*")
	s.Require().NoError(err)
	tmpFile.Close()

	// Get current user ID and group ID for testing
	currentUser, err := user.Current()
	s.Require().NoError(err)

	uid, err := strconv.Atoi(currentUser.Uid)
	s.Require().NoError(err)

	gid, err := strconv.Atoi(currentUser.Gid)
	s.Require().NoError(err)

	// Perform the recursive chown
	err = c.ChownRecursively(tmpDir, uid, gid)
	s.Require().NoError(err)

	// Verify ownership of the directory and file
	verifyOwnership := func(name string) {
		info, err := os.Stat(name)
		s.Require().NoError(err)

		stat, ok := info.Sys().(*syscall.Stat_t)
		s.Require().True(ok, "Expected file info to be syscall.Stat_t")

		s.Equal(uint32(uid), stat.Uid, "UID should match")
		s.Equal(uint32(gid), stat.Gid, "GID should match")
	}

	verifyOwnership(subDir)
	verifyOwnership(tmpFile.Name())
}

func TestFileSystemSuite(t *testing.T) {
	suite.Run(t, new(filesystemSuite))
}
