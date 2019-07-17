package main

import (
	"log"
	"time"
)

func main() {

	msg := time.NewTicker(1*time.Second + 10*time.Millisecond)
	t2 := time.After(1 * time.Second)
	t3 := time.After(2 * time.Second)

	log.Print("start")
outer:
	for {
		select {
		case <-msg.C:
			log.Print("success")
			break outer

		case <-t2:
			log.Print("rx1 timeout")
			continue

		case <-t3:
			log.Print("rx2 timeout")
			continue
		}
	}
	log.Print("end")
}
