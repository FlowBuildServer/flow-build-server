package main

import (
	"context"
	"errors"
	"github.com/google/go-github/github"
	"strings"
)

type Puller struct {
	RepoLink string
	Github   *Github
}

func (p *Puller) validate() error {
	//check if it is github repo
	if !strings.Contains(p.RepoLink, "github.com") {
		return errors.New("Github only!")
	}

	return nil
}

func (p *Puller) Run() ([]*github.PullRequest, error) {
	//build creds (silly way - need to improve)
	parts := strings.Split(p.RepoLink, "/")
	owner, repo := parts[len(parts)-2], parts[len(parts)-1]
	//fetch pulls
	return fetchRecentPullRequests(p.Github, owner, repo)
}

func fetchRecentPullRequests(gh *Github, owner string, repo string) ([]*github.PullRequest, error) {
	var allPullRequests []*github.PullRequest
	ctx := context.Background()
	client := gh.CreateClient()
	options := &github.PullRequestListOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 50},
	}

	// iterate though pages to fetch'em all
	for {
		pullRequests, resp, error := client.PullRequests.List(
			ctx,
			owner,
			repo,
			options,
		)
		if error != nil {
			return nil, error
		} else {
			allPullRequests = append(allPullRequests, pullRequests...)
			if resp.NextPage == 0 {
				break
			}
			options.ListOptions.Page = resp.NextPage
		}
	}

	return allPullRequests, nil
}
