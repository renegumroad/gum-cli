package projectconfig

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/renegumroad/gum-cli/internal/github"
	"github.com/renegumroad/gum-cli/internal/utils"
)

var (
	GumCLiProject = get(&ProjectConfigArgs{
		Name: "gum-cli",
	})

	defaultOwner     = "renegumroad"
	defaultGithubURL = "https://github.com"
)

type ProjectConfigArgs struct {
	Name  string
	URL   string
	Owner string
}

type ProjectConfig interface {
	LatestVersion() (string, error)
}

type projectConfig struct {
	Name  string
	URL   string
	Owner string
	gh    github.Client
}

func get(args *ProjectConfigArgs) ProjectConfig {
	project, err := New(args)

	if err != nil {
		utils.CheckFatalError(err)
	}

	return project
}

func New(args *ProjectConfigArgs) (ProjectConfig, error) {
	return newWithComponents(
		args,
		github.New(),
	)
}

func newWithComponents(args *ProjectConfigArgs, client github.Client) (ProjectConfig, error) {
	if args.Name == "" && args.URL == "" {
		return nil, errors.Errorf("Project name or URL are required to load configuration")
	}

	if args.URL != "" {
		if !strings.HasPrefix(args.URL, defaultGithubURL) {
			return nil, errors.Errorf("Project URL must be a valid GitHub URL")
		}

		owner, repo, err := extractOwnerAndRepoFromURL(args.URL)
		if err != nil {
			return nil, err
		}

		if args.Name == "" {
			args.Name = repo
		} else if args.Name != repo {
			return nil, errors.Errorf("Project name and URL repo name do not match")
		}

		if args.Owner == "" {
			args.Owner = owner
		} else if args.Owner != owner {
			return nil, errors.Errorf("Project owner and URL owner do not match")
		}
	}

	if args.Owner == "" {
		args.Owner = defaultOwner
	}

	if args.URL == "" {
		args.URL = strings.Join([]string{defaultGithubURL, args.Owner, args.Name}, "/")
	}

	config := &projectConfig{
		Name:  args.Name,
		URL:   args.URL,
		Owner: args.Owner,
		gh:    client,
	}

	return config, nil
}

func (p *projectConfig) LatestVersion() (string, error) {
	release, err := p.gh.Projects().GetLatestRelease(&github.ProjectInfo{
		Owner: p.Owner,
		Repo:  p.Name,
	})
	if err != nil {
		return "", errors.Errorf("Failed to get the latest release information: %s", err)

	}
	return release.GetTagName(), nil
}

func extractOwnerAndRepoFromURL(projectURL string) (string, string, error) {
	parsedURL, err := url.Parse(projectURL)
	if err != nil {
		return "", "", err
	}

	pathSegments := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(pathSegments) < 2 {
		return "", "", errors.New("URL path does not contain owner and repo name")
	}

	owner := pathSegments[0]
	repo := pathSegments[1]
	return owner, repo, nil
}
