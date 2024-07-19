package github

import (
	"context"

	"github.com/google/go-github/v63/github"
)

type Release = github.RepositoryRelease

type ProjectInfo struct {
	Owner string
	Repo  string
}

type Projects interface {
	GetLatestRelease(info *ProjectInfo) (*Release, error)
}

type projects struct {
	gh *github.Client
}

func newProjects(gh *github.Client) Projects {
	return &projects{
		gh: gh,
	}
}

func (p *projects) GetLatestRelease(info *ProjectInfo) (*Release, error) {
	ctx := context.Background()
	release, _, err := p.gh.Repositories.GetLatestRelease(ctx, info.Owner, info.Repo)
	if err != nil {
		return nil, err
	}
	return release, nil
}
