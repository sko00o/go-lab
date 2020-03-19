package main

import (
	"encoding/json"
	"fmt"
)

type J struct {
	A bool   `json:"a,omitempty"`
	B uint8  `json:"b,omitempty"`
	C string `json:"c,omitempty"`
}

func main() {
	j := J{
		A: false,
		B: 0,
		C: "",
	}
	bs, err := json.MarshalIndent(j, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))
}
