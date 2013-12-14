package main

import (
	"bytes"
	"encoding/json"
	"github.com/manveru/faker"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"sync/atomic"
	"time"
)

const (
	PubC = 5
	SubC = 20
)

var (
	PubNum int64
	SubNum int64
)

func pub(url string) {
	fake, err := faker.New("en")
	if err != nil {
		log.Panic(err)
	}

	var (
		client http.Client
		msg    = make(map[string]interface{})
		buf    bytes.Buffer
		resp   *http.Response
	)

	for {
		msg["n"] = atomic.AddInt64(&PubNum, 1)
		msg["m"] = fake.Sentence(1, true)
		msg["ts"] = time.Now().Format("15:04:05.000")

		buf.Reset()
		err = json.NewEncoder(&buf).Encode(msg)
		if err != nil {
			log.Panic(err)
		}

		resp, err = client.Post(url, "application/json", &buf)
		if err != nil {
			log.Print(err)
			time.Sleep(time.Second)
			continue
		}
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()

		if resp.StatusCode != 201 {
			log.Panic(resp.Status)
		}
	}
}

func sub(url string) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Panic(err)
	}

	client := http.Client{Jar: jar}

	for {
		resp, err := client.Get(url)
		if err != nil {
			log.Print(err)
			time.Sleep(time.Second)
			continue
		}

		var msg map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&msg)
		if err != nil {
			log.Panic(err)
		}
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()

		atomic.AddInt64(&SubNum, 1)
		// log.Print(msg)
	}
}

func main() {
	url := "http://localhost:8081/topic"

	log.Printf("Publishing messages to %s with concurrency %d.", url, PubC)
	for i := 0; i < PubC; i++ {
		go pub(url)
	}

	log.Printf("Subscribing to messages from %s with concurrency %d.", url, SubC)
	for i := 0; i < SubC; i++ {
		go sub(url)
	}

	for {
		log.Printf("Pub: %10d\tSub: %10d", atomic.LoadInt64(&PubNum), atomic.LoadInt64(&SubNum))
		time.Sleep(time.Second)
	}
}
