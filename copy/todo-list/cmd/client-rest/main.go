package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	address := flag.String("server", "http://localhost:8080", "HTTP gateway url, e.g http://localhost:8080")
	flag.Parse()

	res1 := create(*address)
	read(*address, res1.ID)
	update(*address, res1.ID)
	readAll(*address)
	delete(*address, res1.ID)
}

type created struct {
	API string `json:"api,omitempty"`
	ID  string `json:"id,omitempty"`
}

func create(address string) *created {
	t := time.Now().In(time.UTC)
	pfx := t.Format(time.RFC3339Nano)

	resp, err := http.Post(address+"/v1/todo", "application/json", strings.NewReader(fmt.Sprintf(`
	{
		"api": "v1",
		"toDo": {
			"title": "title (%s)",
			"description": "description: (%s)",
			"reminder": "%s"
		}
	}
	`, pfx, pfx, pfx)))
	if err != nil {
		log.Fatalf("failed to call create method: %v", err)
	}

	var body string
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed to create response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("create response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	var res created
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		log.Fatalf("failed to unmarshal JSON response of Create method: %v", err)
	}
	return &res
}

func read(host string, id string) {
	resp, err := http.Get(fmt.Sprintf("%s%s/%s", host, "/v1/todo", id))
	if err != nil {
		log.Fatalf("failed to call Read method: %v", err)
	}

	var body string
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed read Read response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Read response: Code=%d, Body=%s\n\n", resp.StatusCode, body)
}

func update(host string, id string) {
	t := time.Now().In(time.UTC)
	pfx := t.Format(time.RFC3339Nano)

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s/%s", host, "/v1/todo", id),
		strings.NewReader(fmt.Sprintf(`
	{
		"api":"v1",
		"toDo": {
			"title":"title (%s) + updated",
			"description":"description (%s) + updated",
			"reminder":"%s"
		}
	}
`, pfx, pfx, pfx)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to call Update method: %v", err)
	}
	var body string
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed read Update response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Update response: Code=%d, Body=%s\n\n", resp.StatusCode, body)
}

func delete(host string, id string) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s/%s", host, "/v1/todo", id), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to call Delete method: %v", err)
	}
	var body string
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed read Delete response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Delete response: Code=%d, Body=%s\n\n", resp.StatusCode, body)
}

func readAll(host string) {
	resp, err := http.Get(host + "/v1/todo/all")
	if err != nil {
		log.Fatalf("failed to call read all method: %v", err)
	}
	var body string
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed read read all response: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("ReadAll response: Code=%d, Body=%s\n\n", resp.StatusCode, body)
}
