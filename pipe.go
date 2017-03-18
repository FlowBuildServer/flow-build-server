package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"log"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type Pipe struct {
	GitUrl string
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
	puller := Puller{RepoLink: pipe.GitUrl, User: "", Password: ""}
	err := puller.validate()
	if err != nil {
		log.Println("Not a valid repo", err)
	}
	prs, err := puller.Run()
	if err != nil {
		log.Println(err)
	}

	//Launchers launch pipeline itself
	launchers := make(chan int) //return codes
	var wg sync.WaitGroup
	wg.Add(len(prs))
	for _, pr := range prs {
		go func(pr *github.PullRequest) {
			defer wg.Done()
			//do smth
			log.Println("Building ", *pr.Title, *pr.Head.Label, *pr.Base.Label)

			cmd := exec.Command("docker", "run", "-e", fmt.Sprintf("GIT_REPO=%v", pipe.GitUrl), "mvn")

			if err := cmd.Run(); err != nil {
				log.Println("cmd.Start: ", err)
				exitCode := 1
				if exitError, ok := err.(*exec.ExitError); ok {
					ws := exitError.Sys().(syscall.WaitStatus)
					exitCode = ws.ExitStatus()
				}
				launchers <- exitCode
			} else {
				launchers <- 0
			}
		}(pr)
	}
	//Reporting
	go func() {
		for code := range launchers {
			log.Println("Code", code)
		}
	}()
	wg.Wait()
}
