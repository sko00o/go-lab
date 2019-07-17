package main

import (
	"fmt"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	v, _ := mem.VirtualMemory()

	p := host.HostInfo()

	fmt.Printf("Total: %v, Free: %v, UsedPercent: %f%%\n", v.Total, v.Free, v.UsedPercent)

	fmt.Println(v)
}
