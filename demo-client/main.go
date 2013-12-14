package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fairy-project/fairy/common"
	"github.com/manveru/faker"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func get(url string) {
	for {
		resp, err := http.Get(url)
		check(err)

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		check(err)

		var m common.Message
		err = json.Unmarshal(body, &m)
		check(err)

		for _, v := range m {
			fmt.Println(v)
		}
	}
}

func post(url string, name string, message string) {
	for {
		m := common.Message{"name": name, "message": message}
		fmt.Printf("Posted: %v \n", m)

		b, err := json.Marshal(m)
		check(err)

		var buf bytes.Buffer
		buf.Write([]byte(b))

		resp, err := http.Post(url, "application/json", &buf)
		check(err)

		if resp.StatusCode != 201 {
			fmt.Println(resp.Status)
		}
		resp.Body.Close()

		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	url := "http://localhost:8081/topic"
	fake, err := faker.New("en")
	check(err)

	for i := 0; i < 2; i++ {
		go post(url, fake.Name(), fake.Sentence(5, true))
	}

	get(url)
}
