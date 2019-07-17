package main

import (
	"log"
	"net/http"
	"syscall"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var epoller *epoll

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return
	}

	if err := Add(conn); err != nil {
		log.Printf("Failed to add connetion")
		conn.Close()
	}
}

func SetUlimit() error {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		return err
	}
	rLimit.Cur = rLimit.Max
	return syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
}

func main() {
	// Increase resources limitations
	SetUlimit()

	// Start epoll
	var err error
	epoller, err = MkEpoll()
	if err != nil {
		panic(err)
	}
	go Start()

	recordMetrics()

	//runtime.GOMAXPROCS(runtime.NumCPU)

	http.HandleFunc("/", wsHandler)
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("server on localhost:6666")
	if err := http.ListenAndServe("localhost:6666", nil); err != nil {
		log.Fatal(err)
	}
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

func Start() {
	for {
		connections, err := Wait()
		if err != nil {
			log.Printf("Failed to epoll wait %v", err)
			continue
		}

		for _, conn := range connections {
			if conn == nil {
				break
			}
			if msg, _, err := wsutil.ReadClientData(conn); err != nil {
				if err := Remove(conn); err != nil {
					log.Printf("Failed to remove %v", err)
				}
			} else {
				log.Printf("msg: %s", string(msg))
			}
		}
	}
}
