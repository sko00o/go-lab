package main

import (
	"sync"
)

type Master struct {
	foodCh int
}

type Dog struct {
}

func main() {
	var mailbox uint8
	var lock sync.RWMutex

	sendCond := sync.NewCond(&lock)
	recvCond := sync.NewCond(lock.RLocker())

	// sender
	go func() {
		lock.Lock()
		for mailbox == 1 {
			sendCond.Wait()
		}
		mailbox = 1
		lock.Unlock()
		recvCond.Signal()
	}()

	// receiver
	go func() {
		lock.RLock()
		for mailbox == 0 {
			recvCond.Wait()
		}
		mailbox = 0
		lock.RUnlock()
		sendCond.Signal()
	}()

}
