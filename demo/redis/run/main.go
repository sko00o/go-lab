package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"github.com/sko00o/go-lab/demo/redis/redigo"
	r "github.com/sko00o/go-lab/demo/redis"
)

var (
	URL   string
	Count int

	Delay = time.Millisecond * 200
)

type Packet struct {
	data []byte
}

func init() {
	flag.StringVar(&URL, "url", "redis://redis:6379/1", "redis url")
	flag.IntVar(&Count, "ct", 100, "data count")
	flag.Parse()
}

func main() {
	cfg := r.Config{
		URL:         URL,
		MaxIdle:     10,
		IdleTimeout: time.Minute,
	}
	redigo.Setup(cfg)

	p := redigo.RedisPool()
	idx := 0
	ttl := Delay * 2

	var wg sync.WaitGroup
	for data := range gen(Count, 1*time.Millisecond) {
		idx = (idx + 1) % 10
		key := fmt.Sprintf("collect%d", idx)
		lockKey := fmt.Sprintf("lock%d", idx)

		f := func(data []byte) error {
			start := time.Now()

			if err := Put(p, key, data, ttl); err != nil {
				return err
			}

			if locked, err := Lock(p, lockKey, ttl); err != nil || locked {
				return err
			}

			time.Sleep(Delay)

			d, err := Get(p, key)
			if err != nil {
				return errors.Wrap(err, "get error")
			}
			if len(d) == 0 {
				return fmt.Errorf("get nothing, cost: %v", time.Since(start))
			}
			return nil
		}

		wg.Add(1)
		go func(p Packet) {
			defer wg.Done()
			if err := f(p.data); err != nil {
				fmt.Println("Err: ", err)
				return
			}

			fmt.Printf("%x processed\n", p.data)
		}(Packet{data: data})
	}

	wg.Wait()
}

func gen(l int, delay time.Duration) <-chan []byte {
	ch := make(chan []byte)
	go func() {
		defer close(ch)

		for i := 0; i < l; i++ {
			data := make([]byte, 4)
			i, err := rand.Read(data)
			if err != nil {
				panic(err)
			}
			ch <- data[:i]
			time.Sleep(delay)
		}
	}()

	return ch
}

func Put(p *redis.Pool, key string, val []byte, ttl time.Duration) error {
	c := p.Get()
	defer c.Close()

	c.Send("MULTI")
	c.Send("SADD", key, val)
	c.Send("PEXPIRE", key, int64(ttl)/int64(time.Millisecond))
	_, err := c.Do("EXEC")
	if err != nil {
		return errors.Wrap(err, "put data to set error")
	}

	return nil
}

func Lock(p *redis.Pool, key string, ttl time.Duration) (bool, error) {
	c := p.Get()
	defer c.Close()

	_, err := redis.String(c.Do("SET", key, "lock", "PX", int64(ttl)/int64(time.Millisecond), "NX"))
	if err != nil {
		if err == redis.ErrNil {
			// the packet processing is already locked by an other process
			// so there is nothing to do anymore :-)
			return true, nil
		}
		return false, errors.Wrap(err, "acquire deduplicated lock error")
	}

	return false, nil
}

func Get(p *redis.Pool, key string) ([][]byte, error) {
	c := p.Get()
	defer c.Close()

	return redis.ByteSlices(c.Do("SMEMBERS", key))
}
