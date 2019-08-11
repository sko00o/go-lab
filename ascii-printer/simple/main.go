package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	addr   string
	file   string
	speed  int64
	stdout bool
)

func init() {
	flag.StringVar(&addr, "addr", ":8081", "TCP address to listen on")
	flag.StringVar(&file, "file", "../static/short_intro.txt", "Text file containing the ASCII video")
	flag.Int64Var(&speed, "speed", 15, "Play speed")
	flag.BoolVar(&stdout, "stdout", false, "print to stdout")
	flag.Parse()
}

func main() {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("failed to load file %s: %v\n", file, err)
		return
	}

	frames, err := extract(data)
	if err != nil {
		fmt.Printf("failed to extract: %v\n", err)
		return
	}
	fmt.Printf("extracted %d frames from %s\n", len(frames), file)

	srv := NewServer(frames, speed)
	if stdout {
		srv.play(os.Stdout)
	} else {
		fmt.Printf("listening on %s\n", addr)
		log.Fatal(http.ListenAndServe(addr, srv))
	}
}

type frame struct {
	content  []string
	duration time.Duration
}

func extract(data []byte) ([]frame, error) {
	var frames []frame
	lines := strings.Split(string(data), lineSep)
	for i := 0; i+1+frameHeight <= len(lines); i += frameHeight + 1 {
		dStr := strings.TrimSpace(lines[i])
		dInt, err := strconv.ParseInt(dStr, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("parse frame duration error: %v", err)
		}
		frameDuration := time.Duration(dInt)
		content := lines[i+1 : i+1+frameHeight]

		frames = append(frames, frame{
			content:  content,
			duration: frameDuration,
		})
	}
	return frames, nil
}

var (
	frameHeight = 13
	lineSep     = "\n"
)

type server struct {
	frames []frame
	speed  int64
}

func NewServer(frames []frame, speed int64) *server {
	if speed <= 0 {
		speed = 15
	}
	return &server{
		frames: frames,
		speed:  speed,
	}
}

func (h *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.play(w)
}

func (h *server) play(w io.Writer) {
	for _, frame := range h.frames {
		// Clear terminal and move cursor to position (1,1)
		// TODO: support other terminal
		fmt.Fprint(w, "\033[2J\033[1;1H")
		for _, line := range frame.content {
			fmt.Fprintln(w, line)
		}
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		time.Sleep(frame.duration * time.Second / time.Duration(h.speed))
	}
}
