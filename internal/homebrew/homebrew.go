package homebrew

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/renegumroad/gum-cli/internal/cmdexec"
	"github.com/renegumroad/gum-cli/internal/filesystem"
	"github.com/renegumroad/gum-cli/internal/log"
)

type Package struct {
	Name string
	Cask bool
	Link bool
}

type Client interface {
	Install(pkg Package) error
	IsInstalled(pkg Package) bool
	EnsureInstalled(pkg Package) error
	Link(pkg Package) error
}

type client struct {
	fs filesystem.Client
}

func New() Client {
	return newClientWithComponents(filesystem.New())
}

func newClientWithComponents(fs filesystem.Client) *client {
	return &client{
		fs: fs,
	}
}

func (c *client) EnsureInstalled(pkg Package) error {
	if c.IsInstalled(pkg) {
		log.Infof("Brew package %s is already installed", pkg.Name)
		return nil
	}

	if err := c.Install(pkg); err != nil {
		return err
	}

	if pkg.Link {
		if err := c.Link(pkg); err != nil {
			return err
		}
	}

	log.Infof("Brew package %s installed successfully", pkg.Name)

	return nil
}

func (c *client) Install(pkg Package) error {
	if pkg.Name == "" {
		return errors.Errorf("Package name is required")
	}

	log.Debugf("Installing brew package %s", pkg.Name)
	args := []string{"install"}
	if pkg.Cask {
		args = append(args, "--cask")
	}
	args = append(args, pkg.Name)
	return run(args...)
}

func (c *client) IsInstalled(pkg Package) bool {
	prefix := os.Getenv("HOMEBREW_PREFIX")

	var pkgPath string
	if pkg.Cask {
		pkgPath = filepath.Join(prefix, "Caskroom", pkg.Name)
	} else {
		pkgPath = filepath.Join(prefix, "opt", pkg.Name)
	}

	return c.fs.Exists(pkgPath)
}

func (c *client) Link(pkg Package) error {
	log.Debugf("Linking brew package %s", pkg.Name)
	return run("link", pkg.Name, "--force", "--overwrite")
}

func run(args ...string) error {
	cmd := cmdexec.New("brew", args...)

	err := cmd.Run()

	if err != nil {
		return errors.Errorf("brew %s failed:. err: %s stdout: %s stderr: %s", strings.Join(args, " "), err, cmd.Stdout, cmd.Stderr)
	}

	return nil
}
