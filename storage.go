package main

import (
    "github.com/google/go-github/github"
)

func hasChanges(pr1 *github.PullRequest, pr2 *github.PullRequest) bool {
    return *pr1.Head.SHA != *pr2.Head.SHA || *pr1.Base.SHA != *pr2.Base.SHA
}

type Storage struct {
    All map[int]*github.PullRequest
}

func (self Storage) Has(key int) bool {
    _, ok := self.All[key];

    return ok
}

func (self *Storage) Filter(pullRequests []*github.PullRequest) []*github.PullRequest {
    var out []*github.PullRequest

    for _, pr := range pullRequests {
        if self.Has(*pr.Number) != true {
            self.All[*pr.Number] = pr
            out = append(out, pr)
        } else {
            if hasChanges(self.All[*pr.Number], pr) == true {
                self.All[*pr.Number] = pr
                out = append(out, pr)
            }
        }
    }

    return out
}
