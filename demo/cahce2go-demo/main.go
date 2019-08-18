package main

import (
	"fmt"
	"time"

	"github.com/muesli/cache2go"
)

type myStruct struct {
	text     string
	moreData []byte
}

var (
	somekey        = "somekey"
	expireTime     = time.Second * 3
	getValInterval = time.Second
)

func main() {
	cache := cache2go.Cache("myCache")

	val := myStruct{"Hello", []byte{}}

	cache.Add(somekey, expireTime, &val)

	go func() {
		times := 3
		for {
			select {
			case <-time.After(getValInterval):
				if cache.Exists(somekey) {
					if times > 0 {
						cache.Value(somekey)
						times--
					}
					fmt.Println("check exist")
				} else {
					fmt.Println("check not exist")
					return
				}
			}
		}
	}()

	for {
		select {
		case <-time.After(expireTime):
			if i, err := cache.Value(somekey); err == nil {
				fmt.Println("exist, item:", i)
			} else {
				fmt.Println("not exist, err:", err)
				return
			}
		}
	}
}
