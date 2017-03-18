package main

import(
    "context"
    "github.com/google/go-github/github"
)

type GithubReporter struct {
    Github *Github
}

func (self *GithubReporter) ReportPending(report *Report) (error) {
    return self.setGithubStatus("pending", report)
}

func (self *GithubReporter) ReportSuccess(report *Report) (error) {
    return self.setGithubStatus("success", report)
}

func (self *GithubReporter) ReportError(report *Report) (error) {
    return self.setGithubStatus("error", report)
}

func (self *GithubReporter) setGithubStatus(status string, report *Report) (error) {
    ctx := "flow-build-server"
    message := status
    if report != nil {
        message = report.Message
    }

    _, _, error := self.Github.CreateClient().Repositories.CreateStatus(
        context.Background(),
        *report.PullRequest.Base.Repo.Owner.Login,
        *report.PullRequest.Base.Repo.Name,
        *report.PullRequest.Head.SHA,
        &github.RepoStatus{State: &status, Context: &ctx, Description: &message},
    )

    return error;
}