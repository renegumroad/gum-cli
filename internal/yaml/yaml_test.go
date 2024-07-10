package yaml

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type yamlSuite struct {
	suite.Suite
}

func (s *yamlSuite) TestLoadValidYAML() {
	type Config struct {
		Name string `yaml:"name"`
		Age  int    `yaml:"age"`
	}
	yamlData := []byte("name: John Doe\nage: 30")
	expectedConfig := Config{Name: "John Doe", Age: 30}

	var config Config
	c := New()
	err := c.Load(yamlData, &config)

	s.Require().NoError(err)
	s.Require().Equal(expectedConfig, config)
}

func (s *yamlSuite) TestLoadInvalidYAML() {
	yamlData := []byte("name: John Doe\nage: thirty")

	var config struct{}
	c := New()
	err := c.Load(yamlData, &config)

	s.Require().Error(err)
}

func (s *yamlSuite) TestLoadNilPointer() {
	yamlData := []byte("name: John Doe\nage: 30")
	c := New()
	err := c.Load(yamlData, nil)
	s.Require().Error(err)
}

func (s *yamlSuite) TestReadValidYAMLFile() {
	tmpFile, err := os.CreateTemp("", "valid*.yaml")
	s.Require().NoError(err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("name: John Doe\nage: 30")
	s.Require().NoError(err)
	tmpFile.Close()

	var config struct {
		Name string `yaml:"name"`
		Age  int    `yaml:"age"`
	}
	c := New()
	err = c.Read(tmpFile.Name(), &config)
	s.Require().NoError(err)
	s.Require().Equal("John Doe", config.Name)
	s.Require().Equal(30, config.Age)
}

func (s *yamlSuite) TestReadInvalidYAMLFile() {
	tmpFile, err := os.CreateTemp("", "invalid*.yaml")
	s.Require().NoError(err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("name: John Doe\nage: thirty")
	s.Require().NoError(err)
	tmpFile.Close()

	var config struct{}
	c := New()
	err = c.Read(tmpFile.Name(), &config)
	s.Require().Error(err)
}

func (s *yamlSuite) TestReadNonExistentFile() {
	c := New()
	err := c.Read("nonexistent.yaml", &struct{}{})
	s.Require().Error(err)
}

func (s *yamlSuite) TestReadNilPointer() {
	tmpFile, err := os.CreateTemp("", "valid*.yaml")
	s.Require().NoError(err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("name: John Doe\nage: 30")
	s.Require().NoError(err)
	tmpFile.Close()

	c := New()
	err = c.Read(tmpFile.Name(), nil)
	s.Require().Error(err)
}

func (s *yamlSuite) TestReadDir() {
	tmpDir, err := os.MkdirTemp("", "gum*")
	s.Require().NoError(err)
	defer os.RemoveAll(tmpDir)

	c := New()
	err = c.Read(tmpDir, nil)
	s.Require().Error(err)
}

func TestYamlSuite(t *testing.T) {
	suite.Run(t, new(yamlSuite))
}
