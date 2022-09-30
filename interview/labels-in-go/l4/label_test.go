package main

import (
	"fmt"
)

var println = fmt.Println

func Example_switch() {
	x := 1

	println("before switch")
Switch:
	switch x {
	case 1:
		println("begin case 1")
	Loop:
		for i := 0; i < 3; i++ {
			println("begin for loop")
			switch i {
			case 0:
				println("inner case 0")
				// fallthrough
			case 1:
				println("inner case 1")
				continue Loop
			case 2:
				println("inner case 2")
				break Switch
			}
			println("end for loop")
		}
		println("end case 1")
	}

	println("after switch")
	// Output:
	// before switch
	// begin case 1
	// begin for loop
	// inner case 0
	// end for loop
	// begin for loop
	// inner case 1
	// begin for loop
	// inner case 2
	// after switch
}

func Example_switch_fallthrough() {
	x := 1

	println("before switch")
Switch:
	switch x {
	case 1:
		println("begin case 1")
	Loop:
		for i := 0; i < 3; i++ {
			println("begin for loop")
			switch i {
			case 0:
				println("inner case 0")
				fallthrough // !
			case 1:
				println("inner case 1")
				continue Loop
			case 2:
				println("inner case 2")
				break Switch
			}
			println("end for loop")
		}
		println("end case 1")
	}

	println("after switch")
	// Output:
	// before switch
	// begin case 1
	// begin for loop
	// inner case 0
	// inner case 1
	// begin for loop
	// inner case 1
	// begin for loop
	// inner case 2
	// after switch
}
