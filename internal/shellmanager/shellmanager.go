package shellmanager

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/renehernandez/gum-cli/internal/utils/filesystem"
)

type ShellType string

var (
	ShellZsh  ShellType = "/bin/zsh"
	ShellBash ShellType = "/bin/bash"
	ShellSh   ShellType = "/bin/sh"
)

type Client interface {
	GetShell() string
	GetShellProfilePath(shell ShellType) string
	ProfileByShell() map[ShellType]string
	ProfileContains(shell ShellType, entry string) (bool, error)
	UpdateShellProfile(shell ShellType, entry string) error
}

type client struct {
	fs filesystem.Client
}

func New() Client {
	return newWithComponents(filesystem.NewClient())
}

func newWithComponents(fs filesystem.Client) Client {
	return &client{
		fs: fs,
	}
}

func (c *client) GetShell() string {
	if shell := os.Getenv("SHELL"); shell != "" {
		return shell
	}

	return "/bin/sh"
}

func (c *client) ProfileByShell() map[ShellType]string {
	return map[ShellType]string{
		ShellZsh:  ".zprofile",
		ShellBash: ".bash_profile",
		ShellSh:   ".profile",
	}
}

func (c *client) GetShellProfilePath(shell ShellType) string {
	profile := c.ProfileByShell()[shell]

	homeDir, err := c.fs.HomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(homeDir, profile)
}

func (c *client) UpdateShellProfile(shell ShellType, entry string) error {
	profilePath := c.GetShellProfilePath(shell)

	if err := c.fs.AppendString(profilePath, entry); err != nil {
		return err
	}

	return nil
}

func (c *client) ProfileContains(shell ShellType, entry string) (bool, error) {
	profilePath := c.GetShellProfilePath(shell)

	content, err := os.ReadFile(profilePath)
	if err != nil {
		return false, nil
	}

	return strings.Contains(string(content), entry), nil
}
