package main

import (
	"fmt"
	"log"
	"net/http"
)

type Topic struct {
	channel []interface{}
}

func handler(rw http.ResponseWriter, req *http.Request) {
	topic := req.URL.Path // get topic name
	var topics map[string]Topic
	switch req.Method {
	case "GET":
		save_topic(rw, topic)
	case "POST":
		rw.WriteHeader(201)
	}
	// fmt.Fprintf(rw, req.Method)
}

func save_topic(rw http.ResponseWriter, topic string) {
	fmt.Fprintf(rw, topic)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
