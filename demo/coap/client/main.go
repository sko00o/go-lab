package main

import (
	"context"
	"log"
	"time"

	"github.com/go-ocf/go-coap"
)

func main() {
	co, err := coap.Dial("udp", "0.0.0.0:5688")
	if err != nil {
		log.Fatal("Err dialog: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := co.GetWithContext(ctx, "/a")
	if err != nil {
		log.Fatalf("Err sending request: %v", err)
	}

	log.Printf("Response payload: %s", resp.Payload())
}
