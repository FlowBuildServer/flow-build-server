package main

import (
	"github.com/google/go-github/github"
	"log"
	"sync"
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
	puller := Puller{RepoLink: pipe.GitUrl}
	err := puller.validate()
	if err != nil {
		log.Println("Not a valid repo", err)
	}
	prs, err := puller.run()
	if err != nil {
		log.Println(err)
	}

	//Launchers
	launchers := make(chan int) //return codes
	var wg sync.WaitGroup
	wg.Add(len(prs))
	for _, pr := range prs {
		go func(pr *github.PullRequest) {
			defer wg.Done()
			//do smth
			log.Println("Building ", *pr.Title, *pr.Head.Label, *pr.Base.Label)
			launchers <- 0
		}(pr)
	}
	go func() {
		for code := range launchers {
			log.Println("Success", code)
		}
	}()
	wg.Wait()
	//Reporter
}
