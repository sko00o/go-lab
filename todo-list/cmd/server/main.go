package main

import (
	"fmt"
	"os"

	"github.com/sko00o/go-lab/todo-list/pkg/cmd"
)

func main() {
	if err := cmd.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
