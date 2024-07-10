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
	c := New()
	path, err := c.CurrentDir()

	s.Require().NoError(err)
	s.Require().NotEmpty(path)
}

func (s *filesystemSuite) TestExists() {
	c := New()
	path, _ := c.CurrentDir()

	path = filepath.Join(path, "filesystem.go")

	s.Require().FileExists(path)
}

func (s *filesystemSuite) TestCreateTempDir() {
	c := New()
	tempDir, err := c.CreateTempDir()

	s.Require().NoError(err, "Expected no error on creating a temp dir")
	s.Require().DirExists(tempDir)

	// Clean up after the test
	defer os.RemoveAll(tempDir)

	// Check that the directory name matches the expected pattern
	s.Require().Regexp(`gum_.*`, filepath.Base(tempDir))
}

func (s *filesystemSuite) TestCopyFile() {
	c := New()
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
	c := New()
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
	c := New()
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
	c := New()
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

func (s *filesystemSuite) TestChownDirectory() {
	c := New()
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
	err = c.Chown(tmpDir, uid, gid)
	s.Require().NoError(err)

	// Verify ownership of the directory and file

	s.verifyOwnership(subDir, uid, gid)
	s.verifyOwnership(tmpFile.Name(), uid, gid)
}

func (s *filesystemSuite) TestChownFile() {
	c := New()
	tmpDir, err := c.CreateTempDir()
	s.Require().NoError(err)
	defer os.RemoveAll(tmpDir)

	// Create a subdirectory and a file inside tmpDir to test recursive chown
	tempfile := filepath.Join(tmpDir, "testfile.txt")
	err = c.WriteString(tempfile, "test")
	s.Require().NoError(err)

	// Get current user ID and group ID for testing
	currentUser, err := user.Current()
	s.Require().NoError(err)

	uid, err := strconv.Atoi(currentUser.Uid)
	s.Require().NoError(err)

	gid, err := strconv.Atoi(currentUser.Gid)
	s.Require().NoError(err)

	// Perform the recursive chown
	err = c.Chown(tmpDir, uid, gid)
	s.Require().NoError(err)

	s.verifyOwnership(tempfile, uid, gid)
}

func (s *filesystemSuite) verifyOwnership(name string, uid, gid int) {
	info, err := os.Stat(name)
	s.Require().NoError(err)

	stat, ok := info.Sys().(*syscall.Stat_t)
	s.Require().True(ok, "Expected file info to be syscall.Stat_t")

	s.Equal(uint32(uid), stat.Uid, "UID should match")
	s.Equal(uint32(gid), stat.Gid, "GID should match")
}

func (s *filesystemSuite) TestAppendStringNewFile() {
	c := New()
	tmpDir, err := c.CreateTempDir()
	s.Require().NoError(err)
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "newfile.txt")
	content := "Hello, World!"

	err = c.AppendString(filePath, content)
	s.Require().NoError(err, "Expected no error on appending to a new file")

	readContent, err := os.ReadFile(filePath)
	s.Require().NoError(err, "Expected no error on reading the file")
	s.Require().Equal(content, string(readContent), "Content appended and read should match")
}

func (s *filesystemSuite) TestAppendStringExistingFile() {
	c := New()
	tmpDir, err := c.CreateTempDir()
	s.Require().NoError(err)
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "existingfile.txt")
	initialContent := "Initial Content. "
	additionalContent := "Appended Content."

	// Create file with initial content
	err = os.WriteFile(filePath, []byte(initialContent), 0644)
	s.Require().NoError(err, "Expected no error on creating a file with initial content")

	// Append additional content
	err = c.AppendString(filePath, additionalContent)
	s.Require().NoError(err, "Expected no error on appending to an existing file")

	readContent, err := os.ReadFile(filePath)
	s.Require().NoError(err, "Expected no error on reading the file")
	expectedContent := initialContent + additionalContent
	s.Require().Equal(expectedContent, string(readContent), "Content should include both initial and appended content")
}

func (s *filesystemSuite) TestAppendStringErrorOnInvalidPath() {
	c := New()
	invalidPath := "/invalid/path/to/file.txt"
	err := c.AppendString(invalidPath, "Some content")
	s.Require().Error(err, "Expected an error on appending to a file with an invalid path")
}

func (s *filesystemSuite) TestIsExecutableWithExecutableFile() {
	c := New()
	tmpDir, err := c.CreateTempDir()
	s.Require().NoError(err)
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "executable")
	err = os.WriteFile(filePath, []byte("#!/bin/bash\necho Hello"), 0755)
	s.Require().NoError(err)

	s.Require().True(c.IsExecutable(filePath))
}

func (s *filesystemSuite) TestIsExecutableWithNonExecutableFile() {
	c := New()
	tmpDir, err := c.CreateTempDir()
	s.Require().NoError(err)
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "non_executable")
	err = os.WriteFile(filePath, []byte("Hello, World!"), 0644)
	s.Require().NoError(err)

	s.Require().False(c.IsExecutable(filePath))
}

func (s *filesystemSuite) TestIsExecutableWithNonExistingFile() {
	c := New()
	filePath := "/path/to/non/existing/file"

	s.Require().False(c.IsExecutable(filePath))
}

func (s *filesystemSuite) TestIsExecutableWithDirectory() {
	c := New()
	tmpDir, err := c.CreateTempDir()
	s.Require().NoError(err)
	defer os.RemoveAll(tmpDir)

	s.Require().True(c.IsExecutable(tmpDir))
}

func TestFileSystemSuite(t *testing.T) {
	suite.Run(t, new(filesystemSuite))
}
