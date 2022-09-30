package main

import (
	"fmt"
	"net/http"
)

// error handling in multi goroutine

func main() {

	done := make(chan interface{})
	defer close(done)

	type Result struct {
		Error error
		Resp  *http.Response
	}
	checkState := func(done chan interface{}, urls ...string) <-chan Result {
		res := make(chan Result)

		go func() {
			defer close(res)
			for _, url := range urls {
				var r Result
				r.Resp, r.Error = http.Get(url)
				select {
				case <-done:
					return
				case res <- r:
				}
			}
		}()

		return res
	}

	for res := range checkState(done,
		"https://google.com",
		"https://baidu.com",
		"https://zhihu.com",
		"https://qq.com",
	) {
		if res.Error != nil {
			fmt.Println(res.Error)
			continue
		}
		fmt.Println(res.Resp.Header, res.Resp.Proto)
	}

}
