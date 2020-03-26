package multi_worker_decode

import (
	"bytes"
	"encoding/json"
)

type Simple struct {
	A string
}

type Simple1 struct {
	C Simple
	B int
}

type Simple2 struct {
	B Simple1
	C *Simple1
}

type Complex struct {
	A Simple2
	G string
	H []int
}

func decode(data []byte, v interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	return dec.Decode(&v)
}

func prepare(n int, out chan Complex) func(in []byte) {
	process := make(chan []byte, n)
	for i := 0; i < n; i++ {
		go func() {
			for d := range process {
				var c Complex
				_ = decode(d, &c)
				out <- c
			}
		}()
	}

	return func(in []byte) {
		process <- in
	}
}
