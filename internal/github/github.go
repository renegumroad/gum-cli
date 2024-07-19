package github

import "github.com/google/go-github/v63/github"

type Client interface {
	Projects() Projects
}

type client struct {
	gh       *github.Client
	projects Projects
}

func New() Client {
	gh := github.NewClient(nil)
	c := &client{
		gh:       gh,
		projects: newProjects(gh),
	}

	return c
}

func (c *client) Projects() Projects {
	return c.projects
}
