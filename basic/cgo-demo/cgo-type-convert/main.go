package main

// #include "t.h"
import "C"
import (
	"fmt"
)

func main() {
	b := []byte{1, 2, 3}
	r := test2(b)
	fmt.Println(r)
}

func test1(b []byte) int {
	goUnsafePointer := C.CBytes(b)
	cPointer := (*C.uint8_t)(goUnsafePointer)
	rc := C.test1(cPointer)
	return int(rc)
}

func test2(b []byte) int {
	goUnsafePointer := C.CBytes(b)
	cPointer := (*C.uint8_t)(goUnsafePointer)
	rc := C.test2(cPointer, C.int(len(b)))
	return int(rc)
}
