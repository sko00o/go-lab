package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	player2 "github.com/sko00o/go-lab/product/ascii-printer/player"
	video2 "github.com/sko00o/go-lab/product/ascii-printer/video"
)

var (
	addr   string
	file   string
	stdout bool
)

func init() {
	flag.StringVar(&addr, "addr", ":8080", "TCP address to listen on")
	flag.StringVar(&file, "file", "static/rick_roll.txt", "Text file containing the ASCII video")
	flag.BoolVar(&stdout, "stdout", false, "print to stdout")
	flag.Parse()
}

func main() {
	v := video2.DefaultVideo()
	if err := v.Load(file); err != nil {
		log.Printf("load video error: %v", err)
		return
	}

	if stdout {
		p := player2.DefaultPlayer(os.Stdout, v)
		p.Play()
	} else {
		fmt.Printf("listening on %s\n", addr)
		log.Fatal(http.ListenAndServe(addr, newServer(v)))
	}
}

type server struct {
	v *video2.Video
}

func newServer(v *video2.Video) *server {
	return &server{v}
}

func (h *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := player2.DefaultPlayer(w, h.v)
	p.Play()
}
