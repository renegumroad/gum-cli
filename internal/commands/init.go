package commands

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/renehernandez/gum-cli/internal/log"
	"github.com/renehernandez/gum-cli/internal/utils/filesystem"
	"github.com/renehernandez/gum-cli/internal/version"
)

type InitImpl struct {
	gumHomePath    string
	gumroadPath    string
	gumroadBinPath string
}

func NewInitCmd() *InitImpl {
	return &InitImpl{
		gumroadPath: filepath.Join(filesystem.RootDir(), "opt", "gumroad"),
	}
}

func (cmd *InitImpl) Validate() error {
	if !filesystem.Exists(cmd.gumroadPath) && version.IsRelease() && os.Getuid() != 0 {
		return errors.Errorf("Please re-run gum init with sudo (sudo gum init) to create the necessary directories (/opt/gumroad, /opt/gumroad/bin)")
	}

	return nil
}

func (cmd *InitImpl) Run() error {
	if err := cmd.createGumHomeDir(); err != nil {
		return err
	}

	if err := cmd.createOptGumroadDirs(); err != nil {
		return err
	}

	if err := cmd.CopyExecutable(); err != nil {
		return err
	}

	return nil
}

func (cmd *InitImpl) createGumHomeDir() error {
	if homeDir, err := os.UserHomeDir(); err != nil {
		return errors.Errorf("Unable to get home directory: %s", err)
	} else {
		gumHomePath := filepath.Join(homeDir, ".gum")

		if !filesystem.Exists(gumHomePath) {
			log.Debugln("Creating .gum directory in home directory")
			if err := os.Mkdir(gumHomePath, os.ModePerm); err != nil {
				return errors.Errorf("Unable to create .gum directory in home directory: %s", err)
			}
		}
	}

	return nil
}

func (cmd *InitImpl) createOptGumroadDirs() error {
	if !version.IsRelease() {
		log.Debugln("Skipping creating /opt/gumroad directories since this is not a release version")
		return nil
	}

	if !filesystem.Exists(cmd.gumroadPath) {
		log.Debugf("Creating %s directory", cmd.gumroadPath)
		if err := os.Mkdir(cmd.gumroadPath, os.ModePerm); err != nil {
			return errors.Errorf("Unable to create %s directory: %s", cmd.gumroadPath, err)
		}

		userId := os.Getuid()
		// gid of -1 means that the group id is not changed
		if err := os.Chown(cmd.gumroadPath, userId, -1); err != nil {
			return errors.Errorf("Unable to change ownership of %s directory: %s", cmd.gumroadPath, err)
		}
	}

	cmd.gumroadBinPath = filepath.Join(cmd.gumroadPath, "bin")

	if !filesystem.Exists(cmd.gumroadBinPath) {
		log.Debugf("Creating %s directory", cmd.gumroadBinPath)
		if err := os.Mkdir(cmd.gumroadBinPath, os.ModePerm); err != nil {
			return errors.Errorf("Unable to create %s directory: %s", cmd.gumroadBinPath, err)
		}
	}

	return nil
}

func (cmd *InitImpl) CopyExecutable() error {
	if !version.IsRelease() {
		log.Debugln("Skipping copying gum binary to gumroad bin directory since this is not a release version")
		return nil
	}

	cliPath := filepath.Join(cmd.gumroadBinPath, "gum")

	if !filesystem.Exists(cliPath) {
		log.Debugf("Copying gum binary to %s", cliPath)
		if execPath, err := os.Executable(); err != nil {
			return errors.Errorf("Unable to get current executable path: %s", err)
		} else {
			filesystem.CopyFile(execPath, cliPath)
		}
	}

	return nil
}
