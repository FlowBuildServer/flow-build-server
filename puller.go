package main

import (
    "context"
    "github.com/google/go-github/github"
    "strings"
    "errors"
    //"fmt"
)

type Puller struct {
    RepoLink string
    User string
    Password string
}

type Creds struct {
    Owner string
    Repo string
}

func (p *Puller) validate() (error) {
    //check if it is github repo
    if !strings.Contains(p.RepoLink, "github.com") {
        return errors.New("Github only!")
    }

    return nil
}

func (p *Puller) Run() ([]*github.PullRequest, error) {
    //build creds (silly way - need to improve)
    parts := strings.Split(p.RepoLink, "/")
    owner, repo := parts[len(parts) - 2], parts[len(parts) - 1]
    //fetch pulls
    return fetchRecentPullRequests(p.createClient(), owner, repo)
}

func (p *Puller) createClient() (*github.Client) {
    var client *github.Client
    if p.User != "" && p.Password != "" {
        transport:= github.BasicAuthTransport{p.User, p.Password, "", nil}
        client = github.NewClient(transport.Client())
    } else {
        client = github.NewClient(nil)
    }

    return client
}

func fetchRecentPullRequests(client *github.Client, owner string, repo string) ([]*github.PullRequest, error) {
    var allPullRequests []*github.PullRequest
    ctx := context.Background()
    options := &github.PullRequestListOptions{
        State: "open",
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
