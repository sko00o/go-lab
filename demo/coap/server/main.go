package main

import (
	"context"
	"log"
	"time"

	"github.com/go-ocf/go-coap"
)

func main() {
	mux := coap.NewServeMux()
	log.Println("server running")
	mux.Handle("/a", coap.HandlerFunc(func(w coap.ResponseWriter, r *coap.Request) {
		log.Printf("Got: path=%q: %#v from %v", r.Msg.Path(), r.Msg, r.Client.RemoteAddr())

		w.SetContentFormat(coap.TextPlain)
		log.Printf("Transmitting from A")
		ctx, cancel := context.WithTimeout(r.Ctx, time.Second)
		defer cancel()
		if _, err := w.WriteWithContext(ctx, []byte("hello world")); err != nil {
			log.Printf("Cannot send response: %v", err)
		}
	}))

	log.Fatal(coap.ListenAndServe("udp", ":5688", mux))
}
