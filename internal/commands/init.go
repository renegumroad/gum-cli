package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/renehernandez/gum-cli/assets"
	"github.com/renehernandez/gum-cli/internal/log"
	"github.com/renehernandez/gum-cli/internal/utils/filesystem"
	"github.com/renehernandez/gum-cli/internal/utils/systeminfo"
	"github.com/renehernandez/gum-cli/internal/version"
)

type InitImpl struct {
	homeDir        string
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
	return cmd
}

func (cmd *InitImpl) Validate() error {
	if !cmd.fs.Exists(cmd.gumroadPath) && version.IsRelease() && os.Getuid() != 0 {
		return errors.Errorf("Please re-run gum init with sudo (sudo gum init) to create the necessary directories (/opt/gumroad, /opt/gumroad/bin)")
	}

	return nil
}

func (cmd *InitImpl) Run() error {
	if err := cmd.setup(); err != nil {
		return err
	}

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

	if err := cmd.configureShell(); err != nil {
		return err
	}

	log.Infoln("gum init completed successfully")
	return nil
}

func (cmd *InitImpl) setup() error {
	if homeDir, err := os.UserHomeDir(); err != nil {
		return errors.Errorf("Unable to get home directory: %s", err)
	} else {
		cmd.homeDir = homeDir
		cmd.gumHomePath = filepath.Join(homeDir, ".gum")
		cmd.gumroadPath = filepath.Join(cmd.fs.RootDir(), "opt", "gumroad")
		cmd.gumroadBinPath = filepath.Join(cmd.gumroadPath, "bin")
		cmd.gumroadCliPath = filepath.Join(cmd.gumroadBinPath, "gum")
	}

	return nil
}

func (cmd *InitImpl) createGumHomeDir() error {
	if !cmd.fs.Exists(cmd.gumHomePath) {
		log.Debugln("Creating .gum directory in home directory")
		if err := os.Mkdir(cmd.gumHomePath, os.ModePerm); err != nil {
			return errors.Errorf("Unable to create .gum directory in home directory: %s", err)
		}
	}

	return nil
}

func (cmd *InitImpl) createOptGumroadDirs() error {
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
	if !cmd.fs.Exists(cmd.gumroadCliPath) {
		return errors.Errorf("gum binary does not exist at %s", cmd.gumroadCliPath)
	}

	if !cmd.fs.IsExecutable(cmd.gumroadCliPath) {
		log.Debugf("Making %s executable", cmd.gumroadCliPath)
		if err := cmd.fs.MakeExecutable(cmd.gumroadCliPath); err != nil {
			return err
		}
	} else {
		log.Debugf("%s is already executable", cmd.gumroadCliPath)
	}

	return nil
}

func (cmd *InitImpl) configureShell() error {
	shellConfig, err := assets.GetAsset("shell_config.tmpl")
	if err != nil {
		return errors.Errorf("Unable to retrieve internal shell_config for update: %s", err)
	}

	shellConfigPath := filepath.Join(cmd.gumHomePath, ".shell_config")

	if err := cmd.fs.WriteString(shellConfigPath, string(shellConfig)); err != nil {
		return errors.Errorf("Unable to write shell_config to %s: %s", shellConfigPath, err)
	}

	shell := cmd.sys.GetShell()

	sourceShellEntry := fmt.Sprintf("test -f %s && source %s", shellConfigPath, shellConfigPath)

	var profileFileName string

	if shell == systeminfo.ShellZsh {
		profileFileName = ".zprofile"
	} else {
		profileFileName = ".bash_profile"
	}

	log.Debugf("Checking if shell configuration file %s needs to be updated", profileFileName)
	profilePath := filepath.Join(cmd.homeDir, profileFileName)

	content, err := os.ReadFile(profilePath)
	if err != nil {
		log.Warnf("Unable to read %s file %s: %s. Will create it", profileFileName, profilePath, err)
	}

	if !strings.Contains(string(content), sourceShellEntry) {
		log.Debugf("Adding source %s to %s", shellConfigPath, profilePath)
		if err := cmd.fs.AppendString(profilePath, fmt.Sprintf("\n%s\n", sourceShellEntry)); err != nil {
			return errors.Errorf("Unable to append %s to %s: %s", shellConfigPath, profilePath, err)
		}
	} else {
		log.Debugf("Shell configuration file %s is already updated", profilePath)
	}

	return nil
}
