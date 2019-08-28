package main

import (
	"bytes"
	"encoding/gob"
)

// RPCData
type RPCData struct {
	Name string        // name of the function
	Args []interface{} // request or response' body expect error
	Err  string        // error any executing remote server
}

func Encode(data RPCData) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(b []byte) (RPCData, error) {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	var data RPCData
	if err := decoder.Decode(&data); err != nil {
		return RPCData{}, err
	}
	return data, nil
}
