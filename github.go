package main

import (
    "github.com/google/go-github/github"
)

type Github struct {
    Username string
    Password string
}

func (self *Github) CreateClient() (*github.Client) {
    if self.Username != "" && self.Password != "" {
        t := github.BasicAuthTransport{self.Username, self.Password, "", nil}

        return github.NewClient(t.Client())
    } else {
        return github.NewClient(nil)
    }
}
