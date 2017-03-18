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
	GitUrl string
}

type Step struct {
	PullRequest *github.PullRequest
	ExitCode    int
}

func NewPipe(GitUrl string) Pipe {
	ticker := time.NewTicker(time.Second * 15)
	pipe := Pipe{GitUrl: GitUrl}
	go func() {
		for range ticker.C {
			runPipe(&pipe)
		}
	}()
	return pipe
}

func runPipe(pipe *Pipe) {
	//Puller (Trigger)
	//TODO: move to docker containers
	log.Println("T: ", pipe.GitUrl)
	gh := Github{os.Getenv("GH_LOGIN"), os.Getenv("GH_PASSWORD")}
	puller := Puller{RepoLink: pipe.GitUrl, Github: &gh}
	err := puller.validate()
	if err != nil {
		log.Println("Not a valid repo", err)
	}
	reporter := GithubReporter{&gh}
	prs, err := puller.Run()
	if err != nil {
		log.Println(err)
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
			reporter.ReportPending(&Report{pr, "Pending"})
			cmd := exec.Command("docker", "run", "-e", fmt.Sprintf("GIT_REPO=%v", pipe.GitUrl), "mvn")

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
				reporter.ReportSuccess(&Report{step.PullRequest, "Success"})
			default:
				log.Println("Code", step.ExitCode)
				reporter.ReportError(&Report{step.PullRequest, "Failed"})
			}
		}
	}()
	wg.Wait()
}
