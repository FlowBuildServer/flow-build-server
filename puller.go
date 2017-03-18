package main

import (
    "context"
    "github.com/google/go-github/github"
    //"fmt"
)

type Puller struct {
    //FetchRecentPullRequests func()
}

type Creds struct {
    Owner string
    Repo string
}

func (p *Puller) FetchRecentPullRequests(creds *Creds) ([]*github.PullRequest, error) {
    var allPullRequests []*github.PullRequest
    ctx := context.Background()
    client := github.NewClient(nil)
    options := &github.PullRequestListOptions{
        State: "open",
        ListOptions: github.ListOptions{PerPage: 5},
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
        }
        allPullRequests = append(allPullRequests, pullRequests...)
        if resp.NextPage == 0 {
            break
        }
        options.ListOptions.Page = resp.NextPage
    }

    return allPullRequests, nil
}
