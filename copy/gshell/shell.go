package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if err := execInput(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func echo(c chan string) {
	go func() {
		for {
			fmt.Fprint(os.Stdout, <-c)
		}
	}()
}

// ErrNoPath is returned when 'cd' was called without a second argument.
var ErrNoPath = errors.New("path required")

func execInput(input string) error {
	input = strings.TrimSuffix(input, "\n")
	args := strings.Split(input, " ")
	// check build-in command

	c := make(chan string)
	echo(c)
	switch args[0] {
	case "cd":
		if len(args) < 2 {
			return ErrNoPath
		}
		if err := os.Chdir(args[1]); err != nil {
			return err
		}
		return nil
	case "exit":
		os.Exit(0)
	case "ip":
		fmt.Fprintf(os.Stdout, whereIP(args[1]))
		return nil
	case "echo":
		if len(args) > 1 {
			c <- fmt.Sprint(args[1:])
		}
		return nil
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func whereIP(IP string) string {
	API := "http://ip.cn/index.php?ip="
	r, err := http.Get(API + IP)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	c, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(c)
}
