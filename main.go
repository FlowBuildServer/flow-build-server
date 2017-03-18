package main

import (
	"fmt"
	//"github.com/docker/docker/client"
	//"github.com/levigross/grequests"
	//"gopkg.in/yaml.v2"
)

func main() {
	puller := Puller{}
	pulls, error := puller.FetchRecentPullRequests(&Creds{"google", "ExoPlayer"})

	if error != nil {
	    fmt.Println(error)
	} else {
		fmt.Println(len(pulls))
	}
}
