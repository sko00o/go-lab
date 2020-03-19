package main

import (
	"fmt"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	v, _ := mem.VirtualMemory()
	fmt.Println(v)

	p, _ := host.Info()
	fmt.Println(p)
}
