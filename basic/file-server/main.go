package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var (
	address string
	srvPath string

	sig chan os.Signal
)

/*
# simple file server for upload/download

- download: curl http://srv:80
- upload: curl -X POST -F "file=@/path/to/file" http://srv:80/upload

*/

func init() {
	flag.StringVar(&address, "d", ":80", "listen host:port")
	flag.StringVar(&srvPath, "p", "./", "path for file serve")
	flag.Parse()

	sig = make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(srvPath)))
	http.HandleFunc("/upload", uploadFile)

	go func() {
		log.Printf("serving %s on HTTP port: %s\n", srvPath, address)
		if err := http.ListenAndServe(address, nil); err != nil {
			log.Fatal(err)
		}
	}()

	select {
	case s := <-sig:
		log.Printf("signal received: %v\n", s)
	}
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	// FormFile returns the first file for the given key `file`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "upload failed: %v", err)
		return
	}
	defer file.Close()

	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "upload failed: %v", err)
		return
	}

	if err := ioutil.WriteFile(filepath.Join(srvPath, handler.Filename), fileBytes, 0644); err != nil {
		fmt.Fprintf(w, "upload failed: %v", err)
		return
	}

	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
