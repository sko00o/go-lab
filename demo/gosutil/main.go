package main

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	v, _ := mem.VirtualMemory()
	fmt.Println(v)

	p, _ := host.Info()
	fmt.Println(p)
}
