package commands

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/renehernandez/gum-cli/internal/log"
	"github.com/renehernandez/gum-cli/internal/utils/filesystem"
	"github.com/renehernandez/gum-cli/internal/utils/systeminfo"
	"github.com/renehernandez/gum-cli/internal/version"
)

type InitImpl struct {
	gumHomePath    string
	gumroadPath    string
	gumroadBinPath string
	gumroadCliPath string
	fs             filesystem.Client
	sys            systeminfo.Client
}

func NewInitCmd() *InitImpl {
	cmd := &InitImpl{
		fs:  filesystem.NewClient(),
		sys: systeminfo.NewClient(),
	}

	cmd.gumroadPath = filepath.Join(cmd.fs.RootDir(), "opt", "gumroad")
	cmd.gumroadBinPath = filepath.Join(cmd.gumroadPath, "bin")
	cmd.gumroadCliPath = filepath.Join(cmd.gumroadBinPath, "gum")
	return cmd
}

func (cmd *InitImpl) Validate() error {
	if !cmd.fs.Exists(cmd.gumroadPath) && version.IsRelease() && os.Getuid() != 0 {
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

	if err := cmd.copyExecutable(); err != nil {
		return err
	}

	if err := cmd.setOwnership(); err != nil {
		return err
	}

	if err := cmd.makeCliExecutable(); err != nil {
		return err
	}

	log.Infoln("gum init completed successfully")
	return nil
}

func (cmd *InitImpl) createGumHomeDir() error {
	if homeDir, err := os.UserHomeDir(); err != nil {
		return errors.Errorf("Unable to get home directory: %s", err)
	} else {
		gumHomePath := filepath.Join(homeDir, ".gum")

		if !cmd.fs.Exists(gumHomePath) {
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
		log.Debugln("Skip creating /opt/gumroad directories since this is not a release version")
		return nil
	}

	if !cmd.fs.Exists(cmd.gumroadPath) {
		log.Debugf("Creating %s directory", cmd.gumroadPath)
		if err := os.Mkdir(cmd.gumroadPath, os.ModePerm); err != nil {
			return errors.Errorf("Unable to create %s directory: %s", cmd.gumroadPath, err)
		}
	}

	cmd.gumroadBinPath = filepath.Join(cmd.gumroadPath, "bin")

	if !cmd.fs.Exists(cmd.gumroadBinPath) {
		log.Debugf("Creating %s directory", cmd.gumroadBinPath)
		if err := os.Mkdir(cmd.gumroadBinPath, os.ModePerm); err != nil {
			return errors.Errorf("Unable to create %s directory: %s", cmd.gumroadBinPath, err)
		}
	}

	return nil
}

func (cmd *InitImpl) copyExecutable() error {
	if !version.IsRelease() {
		log.Debugln("Skip copying gum binary to gumroad bin directory since this is not a release version")
		return nil
	}

	cliPath := filepath.Join(cmd.gumroadBinPath, "gum")

	if !cmd.fs.Exists(cliPath) {
		log.Debugf("Copying gum binary to %s", cliPath)
		if execPath, err := os.Executable(); err != nil {
			return errors.Errorf("Unable to get current executable path: %s", err)
		} else {
			cmd.fs.CopyFile(execPath, cliPath)
		}
	}

	return nil
}

func (cmd *InitImpl) setOwnership() error {
	if !version.IsRelease() {
		log.Debugln("Skip ownership change since this is not a release version")
		return nil
	}

	log.Debugf("Checking if ownership of %s directory needs to be updated", cmd.gumroadPath)

	info, err := cmd.fs.GetOwner(cmd.gumroadPath)
	if err != nil {
		return err
	}

	if info.Id != 0 {
		log.Debugf("Won't update ownership of %s directory since it is not owned by root. Owner: %s", cmd.gumroadPath, info.Name)
		return nil
	}

	user, err := cmd.sys.GetSudoOriginalUser()
	if err != nil {
		return err
	}
	log.Debugf("Setting ownership of %s directory to %s", cmd.gumroadPath, user.Name)

	return cmd.fs.ChownUserRecursively(cmd.gumroadPath, user.Id)
}

func (cmd *InitImpl) makeCliExecutable() error {
	if !version.IsRelease() {
		log.Debugf("Skip making gum executable at %s since this is not a release version", cmd.gumroadCliPath)
		return nil
	}

	if !cmd.fs.Exists(cmd.gumroadCliPath) {
		return errors.Errorf("gum binary does not exist at %s", cmd.gumroadCliPath)
	}

	log.Debugf("Making %s executable", cmd.gumroadCliPath)
	return cmd.fs.MakeExecutable(cmd.gumroadCliPath)
}
