package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"reflect"
)

// RPCServer is RPC server definition
type RPCServer struct {
	addr  string
	funcs map[string]reflect.Value
}

var logServer = log.New(os.Stderr, "server: ", log.LstdFlags)

func NewServer(addr string) *RPCServer {
	return &RPCServer{
		addr:  addr,
		funcs: make(map[string]reflect.Value),
	}
}

func (s *RPCServer) Register(fName string, f interface{}) {
	if _, ok := s.funcs[fName]; ok {
		return
	}

	s.funcs[fName] = reflect.ValueOf(f)
}

func (s *RPCServer) Execute(req RPCData) RPCData {
	f, ok := s.funcs[req.Name]
	if !ok {
		e := fmt.Sprintf("func %s is not register", req.Name)
		logServer.Println(e)
		return RPCData{Name: req.Name, Err: e}
	}

	logServer.Printf("func %s is called", req.Name)

	inArgs := make([]reflect.Value, len(req.Args))
	for i := range req.Args {
		inArgs[i] = reflect.ValueOf(req.Args[i])
	}

	out := f.Call(inArgs)

	resArgs := make([]interface{}, len(out)-1)
	for i := 0; i < len(out)-1; i++ {
		resArgs[i] = out[i].Interface()
	}

	var er string
	if e, ok := out[len(out)-1].Interface().(error); ok {
		er = e.Error()
	}
	return RPCData{Name: req.Name, Args: resArgs, Err: er}
}

func (s *RPCServer) Run() {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		logServer.Printf("listen on %s error: %s", s.addr, err)
		return
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			logServer.Printf("accept error: %v", err)
			continue
		}

		go func() {
			connTransport := NewTransport(conn)
			for {
				req, err := connTransport.Read()
				if err != nil && err != io.EOF {
					logServer.Printf("read error: %v", err)
					return
				}

				decReq, err := Decode(req)
				if err != nil {
					logServer.Printf("error decode the payload: %v", err)
					return
				}

				resp := s.Execute(decReq)
				b, err := Encode(resp)
				if err != nil {
					logServer.Printf("error encode the payload for response: %v", err)
					return
				}

				err = connTransport.Send(b)
				if err != nil {
					logServer.Printf("transport write error: %v", err)
				}
			}
		}()
	}
}
