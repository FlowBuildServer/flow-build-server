package main

import (
	"log"
	"time"
)

type Pipe struct {
	GitUrl string
}

func NewPipe(GitUrl string) Pipe {
	ticker := time.NewTicker(time.Second * 1)
	pipe := Pipe{GitUrl: GitUrl}
	go func() {
		for range ticker.C {
			runPipe(&pipe)
		}
	}()
	return pipe
}

func runPipe(pipe *Pipe) {
	//Puller
	log.Println("T: ", pipe.GitUrl)
	//Launcher

	//Reporter
}
