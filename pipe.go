package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type Pipe struct {
	GitUrl     string
	PipePuller *Puller
	GithubRepoter *GithubReporter
	TelegramRepoter *TelegramReporter
}

type Step struct {
	PullRequest *github.PullRequest
	ExitCode    int
}

func NewPipe(GitUrl string) (*Pipe, error) {
	ticker := time.NewTicker(time.Second * 10)
	gh := Github{os.Getenv("GH_LOGIN"), os.Getenv("GH_PASSWORD")}
	t := Telegram{"http://localhost:8080"}
	puller := &Puller{RepoLink: GitUrl, Github: &gh, Storage: &Storage{make(map[int]*github.PullRequest)}}
	//validate before run pipe
	err := puller.Validate()
	if err != nil {
		log.Println("Not a valid repo", err)

		return nil, err
	}
	pipe := &Pipe{
		GitUrl: GitUrl,
		PipePuller: puller,
		GithubRepoter: &GithubReporter{&gh},
		TelegramRepoter: &TelegramReporter{&t},
	}
	go func() {
		for range ticker.C {
			runPipe(pipe)
		}
	}()
	return pipe, nil
}

func runPipe(pipe *Pipe) {
	//Puller (Trigger)
	//TODO: move to docker containers
	prs, err := pipe.PipePuller.Run()
	if err != nil {
		log.Println(err)
		return
	}
	//Launchers launch pipeline itself
	launchers := make(chan Step) //return codes
	var wg sync.WaitGroup
	wg.Add(len(prs))
	for _, pr := range prs {
		go func(pr *github.PullRequest) {
			defer wg.Done()
			//do smth
			log.Println("Building ", *pr.Title, *pr.Head.Label, *pr.Base.Label)
			pipe.GithubRepoter.ReportPending(&Report{pr, "Pending"})
			//pipe.TelegramRepoter.ReportPending(&Report{pr, "Pending"})
			cmd := exec.Command("docker", "run", "-e", fmt.Sprintf("GIT_REPO=%v", pipe.GitUrl), "-e", fmt.Sprintf("SOURCE_BRANCH=%v", *pr.Head.Ref), "-e", fmt.Sprintf("TARGET_BRANCH=%v", *pr.Base.Ref), "mvn")

			if err := cmd.Run(); err != nil {
				log.Println("cmd.Start: ", err)
				exitCode := 1
				if exitError, ok := err.(*exec.ExitError); ok {
					ws := exitError.Sys().(syscall.WaitStatus)
					exitCode = ws.ExitStatus()
				}
				launchers <- Step{PullRequest: pr, ExitCode: exitCode}
			} else {
				launchers <- Step{PullRequest: pr, ExitCode: 0}
			}
		}(pr)
	}
	//Reporting
	go func() {
		for step := range launchers {
			switch step.ExitCode {
			case 0:
				log.Println("Code", step.ExitCode)
				pipe.GithubRepoter.ReportSuccess(&Report{step.PullRequest, "Success"})
			default:
				log.Println("Code", step.ExitCode)
				pipe.GithubRepoter.ReportError(&Report{step.PullRequest, "Failed"})
			}
		}
	}()
	wg.Wait()
}
