package main

import (
	"context"
	"errors"
	"github.com/google/go-github/github"
	"strings"
	//"fmt"
)

type Puller struct {
	RepoLink string
}

type Creds struct {
	Owner string
	Repo  string
}

func (p *Puller) validate() error {
	//check if it is github repo
	if !strings.Contains(p.RepoLink, "github.com") {
		return errors.New("Github only!")
	}

	return nil
}

func (p *Puller) run() ([]*github.PullRequest, error) {
	//build creds (silly way - need to improve)
	parts := strings.Split(p.RepoLink, "/")
	creds := Creds{parts[len(parts)-2], parts[len(parts)-1]}
	//fetch pulls
	return fetchRecentPullRequests(&creds)
}

func fetchRecentPullRequests(creds *Creds) ([]*github.PullRequest, error) {
	var allPullRequests []*github.PullRequest
	ctx := context.Background()
	client := github.NewClient(nil)
	options := &github.PullRequestListOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 50},
	}

	// iterate though pages to fetch'em all
	for {
		pullRequests, resp, error := client.PullRequests.List(
			ctx,
			creds.Owner,
			creds.Repo,
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

func fetchRepo(creds *Creds) (*github.Repository, error) {
	ctx := context.Background()
	client := github.NewClient(nil)
	repo, _, error := client.Repositories.Get(
		ctx,
		creds.Owner,
		creds.Repo,
	)
	if error != nil {
		return nil, error
	}

	return repo, nil
}
