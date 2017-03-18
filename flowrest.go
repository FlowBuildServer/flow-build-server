package main

import (
	//"fmt"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	//"github.com/docker/docker/client"
	//"github.com/google/go-github/github"
	//"github.com/levigross/grequests"
	//"gopkg.in/yaml.v2"
)

var pipechans chan int
var pipes []Pipe

type pipe_post struct {
	URL string
}

func flowbuild(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Jenkins killa\n"))
}

func addnewpipe(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var in pipe_post
	err := decoder.Decode(&in)
	if err != nil {
		log.Println("Malformed request", err)
	}
	log.Println(in)
	pipes = append(pipes, NewPipe(in.URL))
}
func waitforpipes() {
	<-pipechans
	log.Println("Closing")
	return
}

func StartFlowBuilder() {
	defer waitforpipes()
	pipes = make([]Pipe, 20)
	r := mux.NewRouter()
	r.HandleFunc("/", flowbuild).Methods("GET")
	r.HandleFunc("/newpipe", addnewpipe).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", r))
}
