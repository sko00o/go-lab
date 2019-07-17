package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	srvPath string
	listen  string
)

func init() {
	flag.StringVar(&srvPath, "s", ".", "serve file path")
	flag.StringVar(&listen, "l", ":8089", "serve listen port")
	flag.Parse()
}

func main() {
	log.Println("serve file on", srvPath, ", listen on", listen)
	if err := http.ListenAndServe(listen, http.FileServer(http.Dir(srvPath))); err != nil {
		log.Fatal(err)
	}
}
