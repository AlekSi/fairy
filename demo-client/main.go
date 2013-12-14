package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
  "github.com/manveru/faker"
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

		var m map[string]interface{}
		err = json.Unmarshal(body, &m)
		check(err)

		for _, v := range m { fmt.Println(v) }
	}
}

func post(url string, name string, message string) {
	for {
		m := map[string]interface{}{"name": name, "message": message}
		fmt.Println(m)

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
	url := "https://rubygems.org/api/v1/gems/rails.json"
  fake, err := faker.New("en")
  check(err)

	for i := 0; i < 2; i++ {
		go post(url, fake.Name(), fake.Sentence(5, true))
	}

	get(url)
}
