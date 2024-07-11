package bundler

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/renegumroad/gum-cli/internal/cli/cmdexec"
	"github.com/renegumroad/gum-cli/internal/filesystem"
	"github.com/renegumroad/gum-cli/internal/log"
)

type Client interface {
	InstallGems() error
	InstallBundler() error
	IsBundlerInstalled() bool
	EnsureBundlerInstalled() error
}

type client struct {
	dir    string
	fs     filesystem.Client
	cmdGen cmdexec.CmdGenerator
}

func New() Client {
	return newClientWithComponents(
		filesystem.New(),
		cmdexec.NewCommandGenerator(),
	)
}

func newClientWithComponents(
	fs filesystem.Client,
	cmdGen cmdexec.CmdGenerator,
) *client {
	return &client{
		fs:     fs,
		cmdGen: cmdGen,
	}
}

func (c *client) InstallGems() error {
	log.Debugf("Installing gems with Bundler")
	dir, err := c.fs.CurrentDir()
	if err != nil {
		return err
	}
	gemfile := filepath.Join(dir, "Gemfile")

	if !c.fs.Exists(gemfile) {
		return errors.Errorf("Gemfile not found in %s", c.dir)
	}

	log.Infof("Running bundle install")
	cmd := c.cmdGen("bundle", "install")

	if err := cmd.Run(); err != nil {
		return errors.Errorf("Failed to install gems: err: %s; stdout: %s; stderr: %s", err, cmd.Stdout(), cmd.Stderr())
	}

	return nil
}

func (c *client) EnsureBundlerInstalled() error {
	if c.IsBundlerInstalled() {
		log.Infof("Bundler is already installed")
		return nil
	}

	if err := c.InstallBundler(); err != nil {
		return err
	}

	return nil
}

func (c *client) IsBundlerInstalled() bool {
	version := c.getBundlerVersion()

	if version == "" {
		log.Debugf("No bundler version found. Bundler will use the default version")
		return true
	}

	cmd := c.cmdGen("gem", "list", "--installed", "--exact", "bundler", "--version", version)

	if err := cmd.Run(); err != nil {
		return false
	}

	return strings.EqualFold(cmd.Stdout(), "true")
}

func (c *client) InstallBundler() error {
	version := c.getBundlerVersion()

	if version == "" {
		log.Debugf("No bundler version found. Using default version")
		return nil
	}

	log.Infof("Installing Bundler version %s", version)

	cmd := c.cmdGen("gem", "install", fmt.Sprintf("bundler:%s", version))

	if err := cmd.Run(); err != nil {
		return errors.Errorf("Failed to install bundler gem: %s", err)
	}

	return nil
}

func (c *client) getBundlerVersion() string {
	version, err := c.getVersionFromVersionFile()

	if err != nil {
		version, _ = c.getVersionFromGemfileLock()
	}

	return version
}

func (c *client) getVersionFromVersionFile() (string, error) {
	dir, err := c.fs.CurrentDir()
	if err != nil {
		return "", err
	}
	bundlerVersionFile := filepath.Join(dir, ".bundler-version")
	log.Debugf("Checking for Bundler version in %s", bundlerVersionFile)

	if !c.fs.Exists(bundlerVersionFile) {
		return "", nil
	}

	content, err := c.fs.ReadString(bundlerVersionFile)
	if err != nil {
		return "", err
	}

	if content == "" {
		return "", errors.Errorf("Bundler version not found in .bundler-version file")
	}

	re := regexp.MustCompile(`^\d+\.\d+\.\d+$`)

	match := re.MatchString(content)

	if !match {
		return "", errors.Errorf("Invalid Bundler version in .bundler-version file")
	}

	return content, nil
}

func (c *client) getVersionFromGemfileLock() (string, error) {
	dir, err := c.fs.CurrentDir()
	if err != nil {
		return "", err
	}
	gemfileLock := filepath.Join(dir, "Gemfile.lock")
	log.Debugf("Checking for Bundler version in %s", gemfileLock)

	if !c.fs.Exists(gemfileLock) {
		return "", nil
	}

	content, err := c.fs.ReadString(gemfileLock)
	if err != nil {
		return "", err
	}

	// Define a regular expression to find the Bundler version
	// The version is expected to be in the format of `BUNDLED WITH` followed by the version number
	re := regexp.MustCompile(`BUNDLED WITH\s+(\d+\.\d+\.\d+)`)

	// Find the version using the regular expression
	matches := re.FindStringSubmatch(content)

	// Check if a version was found
	if len(matches) < 2 {
		return "", errors.Errorf("Bundler version not found in Gemfile.lock")
	}

	return matches[1], nil
}
