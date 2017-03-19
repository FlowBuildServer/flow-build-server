package main

import (
	"fmt"
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
var pipes []*Pipe

type pipe_post struct {
	URL string
}

func flowbuild(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Jenkins killa\n"))
}

func addnewpipe(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var in pipe_post
	err := decoder.Decode(&in)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("{\"message\": \"Malformed request\"}"))
		return
	}

	log.Println(in)
	pipe, err := NewPipe(in.URL)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("{\"message\": \"%v\"}", err)))
		return
	} else {
		pipes = append(pipes, pipe)
		w.WriteHeader(200)
		w.Write([]byte("{\"message\": \"OK\"}"))
	}
}
func waitforpipes() {
	<-pipechans
	log.Println("Closing")
	return
}

func StartFlowBuilder() {
	defer waitforpipes()
	pipes = make([]*Pipe, 20)
	r := mux.NewRouter()
	r.HandleFunc("/", flowbuild).Methods("GET")
	r.HandleFunc("/newpipe", addnewpipe).Methods("POST")
	log.Println("Flow Build Server started!")
	log.Fatal(http.ListenAndServe(":8000", r))
}
