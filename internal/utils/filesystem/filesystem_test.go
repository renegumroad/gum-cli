package filesystem

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type filesystemSuite struct {
	suite.Suite
}

func (s *filesystemSuite) TestCurrentDir() {
	c := NewClient()
	path, err := c.CurrentDir()

	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), path)
}

func (s *filesystemSuite) TestExists() {
	c := NewClient()
	path, _ := c.CurrentDir()

	path = filepath.Join(path, "filesystem.go")

	fmt.Printf("Path: %s\n", path)

	assert.True(s.T(), c.Exists(path))
}

func TestFilesystemSuite(t *testing.T) {
	suite.Run(t, new(filesystemSuite))
}
