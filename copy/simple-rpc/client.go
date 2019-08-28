package main

import (
	"errors"
	"net"
	"reflect"
)

type Client struct {
	conn net.Conn
}

func NewClient(conn net.Conn) *Client {
	return &Client{conn: conn}
}

func (c *Client) callRPC(rpcName string, fPtr interface{}) {
	container := reflect.ValueOf(fPtr).Elem()
	f := func(req []reflect.Value) []reflect.Value {
		cReqTransport := NewTransport(c.conn)

		errorHandler := func(err error) []reflect.Value {
			outArgs := make([]reflect.Value, container.Type().NumOut())
			for i := 0; i < len(outArgs)-1; i++ {
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
			outArgs[len(outArgs)-1] = reflect.ValueOf(&err).Elem()
			return outArgs
		}

		inArgs := make([]interface{}, 0, len(req))
		for _, arg := range req {
			inArgs = append(inArgs, arg.Interface())
		}

		reqRPC := RPCData{Name: rpcName, Args: inArgs}
		b, err := Encode(reqRPC)
		if err != nil {
			panic(err)
		}

		err = cReqTransport.Send(b)
		if err != nil {
			return errorHandler(err)
		}

		rsp, err := cReqTransport.Read()
		if err != nil {
			return errorHandler(err)
		}

		resp, err := Decode(rsp)
		if err != nil {
			panic(err)
		}

		if resp.Err != "" {
			return errorHandler(errors.New(resp.Err))
		}

		numOut := container.Type().NumOut()
		if len(resp.Args) == 0 {
			resp.Args = make([]interface{}, numOut)
		}

		outArgs := make([]reflect.Value, numOut)
		for i := 0; i < numOut; i++ {
			if i != numOut-1 {
				// if argument is nil (gob will ignore "Zero" in transmission), set "Zero" value
				if resp.Args[i] == nil {
					outArgs[i] = reflect.Zero(container.Type().Out(i))
				} else {
					outArgs[i] = reflect.ValueOf(resp.Args[i])
				}
			} else {
				// unpack error argument
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
		}

		return outArgs
	}
	container.Set(reflect.MakeFunc(container.Type(), f))
}
