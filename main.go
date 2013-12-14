package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func handler(rw http.ResponseWriter, req *http.Request) {
	topic := req.URL.Path
  switch req.Method {
    rw.WriteHeader(201)
  }

}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
