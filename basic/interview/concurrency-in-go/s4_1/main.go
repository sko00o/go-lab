package main

import (
	"fmt"
	"time"
)

// 'or' channel

func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	// 相当于每个节点管理最多4个 channel ，这产生一棵树，任意 channel 关闭都可以导致所有 channel 关闭
	orDone := make(chan interface{})
	go func() {
		defer close(orDone)
		switch len(channels) {
		case 2:
			switch {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-or(append(channels[3:], orDone)...):
			}
		}
	}()
	return orDone
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		cc := make(chan interface{})
		go func() {
			time.Sleep(after)
			close(cc)
		}()
		return cc
	}

	start := time.Now()
	<-or(
		sig(time.Hour),
		sig(time.Minute),
		sig(time.Second),
	)

	fmt.Printf("done after: %v", time.Since(start))
}
