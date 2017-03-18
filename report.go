package main

import (
    "github.com/google/go-github/github"
)

type Report struct {
    PullRequest *github.PullRequest
    Message string
}
