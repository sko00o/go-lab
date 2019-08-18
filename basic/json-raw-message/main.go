package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type TestJson struct {
	Object json.RawMessage `json:"object"`
}

func main() {

	task := []struct {
		name string
		in   TestJson
	}{
		{
			"no-obj",
			TestJson{},
		},
		{
			"empty-json-obj",
			TestJson{[]byte("{}")},
		},
	}

	for _, t := range task {
		b, err := json.Marshal(t.in)
		if err != nil {
			fmt.Println("marshal error", err)
		}
		var out TestJson
		if err := json.Unmarshal(b, &out); err != nil {
			fmt.Println("unmarshal error", err)
		}

		if !reflect.DeepEqual(t.in, out) {
			fmt.Printf("task: %s, different, in: %v, out: %v \n", t.name, t.in, out)
		} else {
			fmt.Printf("task: %s, in and out is same", t.name)
		}
	}
}
