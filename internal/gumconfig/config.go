package gumconfig

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/renegumroad/gum-cli/internal/actions"
	"github.com/renegumroad/gum-cli/internal/cli/homebrew"
	"github.com/renegumroad/gum-cli/internal/filesystem"
	"github.com/renegumroad/gum-cli/internal/log"
	"github.com/renegumroad/gum-cli/internal/yaml"
)

var (
	configFileNameOptions = []string{
		"gum.yml",
		"gum.yaml",
	}
)

type GumConfig struct {
	Up []UpAction `yaml:"up,omitempty"`
}

type UpAction struct {
	Action NamedAction        `yaml:"action,omitempty"`
	Brew   []homebrew.Package `yaml:"brew,omitempty"`
}

type NamedAction string

func New(dir string) (*GumConfig, error) {
	config, err := findConfig(dir)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (config *GumConfig) Validate() error {
	log.Debugf("Validating gum config")

	for _, up := range config.Up {
		if up.Action == "" && len(up.Brew) == 0 {
			return errors.Errorf("Named action or brew packages are required")
		} else if up.Action != "" && len(up.Brew) > 0 {
			return errors.Errorf("Cannot defined a named action and brew packages in the same entry")
		}

		if up.Action != "" {
			if !actions.SupportedByConfig(string(up.Action)) {
				return errors.Errorf("Named action %s does not exist or cannot be invoked via config", up.Action)
			}
		}

		for _, pkg := range up.Brew {
			if pkg.Name == "" {
				return errors.Errorf("Package name is required")
			}
		}

	}

	log.Infoln("gum.yml config validated successfully")
	return nil
}

func findConfig(dir string) (*GumConfig, error) {
	log.Debugf("Detecting gum config in %s", dir)
	fs := filesystem.New()
	for _, fileName := range configFileNameOptions {
		path := filepath.Join(dir, fileName)

		if fs.Exists(path) {
			return parseConfig(path)
		}
	}

	return nil, errors.Errorf("No config file found in %s. Expected filenames: %s", dir, configFileNameOptions)
}

func parseConfig(path string) (*GumConfig, error) {
	log.Debugf("Parsing gum config file: %s", path)
	config := &GumConfig{}

	client := yaml.New()

	if err := client.Read(path, config); err != nil {
		return nil, err
	}

	return config, nil
}
