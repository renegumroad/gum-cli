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
	path, err := CurrentDir()

	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), path)
}

func (s *filesystemSuite) TestExists() {
	path, _ := CurrentDir()

	path = filepath.Join(path, "filesystem.go")

	fmt.Printf("Path: %s\n", path)

	assert.True(s.T(), Exists(path))
}

func TestFilesystemSuite(t *testing.T) {
	suite.Run(t, new(filesystemSuite))
}
