package main

import (
	"encoding/json"
	"github.com/fairy-project/fairy/common"
	"log"
	"net/http"
)

var (
	topics = make(map[string]Topic)
)

func publish(rw http.ResponseWriter, req *http.Request) {
	var msg common.Message
	err := json.NewDecoder(req.Body).Decode(&msg)
	req.Body.Close()
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	t := topics[req.URL.Path]
	t.Publish(msg)

	rw.WriteHeader(201)
}

func subscribe(rw http.ResponseWriter, req *http.Request) {
	req.Body.Close()

	t, ok := topics[req.URL.Path]
	if !ok {
		t = NewTopic()
		topics[req.URL.Path] = t
	}

	c := t.GetChannel(req.RemoteAddr)
	msg := <-c

	json.NewEncoder(rw).Encode(msg)
}

func handler(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		subscribe(rw, req)
	case "POST":
		publish(rw, req)
	default:
		http.Error(rw, "bad method", 400)
	}
}

func main() {
	adminPrefix := "/meta/admin/"
	http.Handle(adminPrefix, http.StripPrefix(adminPrefix, http.FileServer(http.Dir("public"))))

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
