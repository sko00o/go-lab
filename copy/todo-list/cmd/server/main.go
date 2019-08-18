package main

import (
	"fmt"
	"os"

	cmd2 "github.com/sko00o/go-lab/copy/todo-list/pkg/cmd"
)

func main() {
	if err := cmd2.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
