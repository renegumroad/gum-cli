package update

import (
	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"github.com/renegumroad/gum-cli/internal/log"
	"github.com/renegumroad/gum-cli/internal/projectconfig"
	"github.com/renegumroad/gum-cli/internal/version"
)

type UpdateImpl struct {
	gumProject projectconfig.ProjectConfig
}

func New() *UpdateImpl {
	return newWithComponents(projectconfig.GumCLiProject)
}

func newWithComponents(projectConfig projectconfig.ProjectConfig) *UpdateImpl {
	return &UpdateImpl{
		gumProject: projectConfig,
	}
}

func (impl *UpdateImpl) Validate() error {
	return nil
}

func (impl *UpdateImpl) Run() error {
	needsUpdate, err := impl.checkforUpdate()
	if err != nil {
		return err
	}

	if !needsUpdate {
		log.Infof("You are already on the latest version of gum cli")
		return nil
	}

	return nil
}

func (impl *UpdateImpl) checkforUpdate() (bool, error) {
	codeSemver, err := semver.NewVersion(version.VERSION)
	if err != nil {
		return false, errors.Errorf("Current version is not a valid semver: %s. Won't attempt to upgrade", version.VERSION)
	}

	latestVersion, err := impl.gumProject.LatestVersion()
	if err != nil {
		return false, err
	}

	if latestVersion == "" {
		return false, nil
	}

	latestSemver, err := semver.NewVersion(latestVersion)
	if err != nil {
		return false, errors.Errorf("Latest version is not a valid semver: %s. Won't attempt to upgrade", latestVersion)
	}

	return codeSemver.LessThan(latestSemver), nil
}
