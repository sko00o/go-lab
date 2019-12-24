package test_select

import (
	"fmt"
	"net/http"
	"time"
)

func RacerOld(a, b string) (winner string) {
	aDuration := measureResponseTime(a)
	bDuration := measureResponseTime(b)

	if aDuration < bDuration {
		return a
	}

	return b
}

func measureResponseTime(url string) time.Duration {
	startB := time.Now()
	http.Get(url)
	bDuration := time.Since(startB)
	return bDuration
}

func Racer(a, b string) (winner string, error error) {
	return ConfigurableRacer(a, b, tenSecondTimeout)
}

var tenSecondTimeout = 10 * time.Second

func ConfigurableRacer(a, b string, delay time.Duration) (winner string, err error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	case <-time.After(delay):
		return "", fmt.Errorf("timeout wait %s and %s", a, b)
	}
}

func ping(url string) chan bool {
	ch := make(chan bool)

	go func() {
		http.Get(url)
		ch <- true
	}()

	return ch
}
