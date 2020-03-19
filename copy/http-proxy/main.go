package main

import (
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

var (
	pemPath string
	keyPath string
	proto   string
	addr    string
)

func init() {
	flag.StringVar(&pemPath, "pem", "server.pem", "path to pem path")
	flag.StringVar(&keyPath, "key", "server.key", "path to key path")
	flag.StringVar(&proto, "proto", "http", "proxy protocol ( http or https )")
	flag.StringVar(&addr, "addr", ":8888", "listen address")
	flag.Parse()
}

func main() {
	server := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				handleTunneling(w, r)
			} else {
				handleHTTP(w, r)
			}
		}),
		TLSNextProto: make(map[string]func(server *http.Server, conn *tls.Conn, handler http.Handler)),
	}

	log.Println("listening on", addr)
	if proto == "http" {
		log.Fatal(server.ListenAndServe())
	} else {
		log.Fatal(server.ListenAndServeTLS(pemPath, keyPath))
	}
}

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	dst, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "hijack not support", http.StatusServiceUnavailable)
		return
	}

	src, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	go transfer(dst, src)
	go transfer(src, dst)
}

func transfer(dst io.WriteCloser, src io.ReadCloser) {
	defer dst.Close()
	defer src.Close()
	io.Copy(dst, src)
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
