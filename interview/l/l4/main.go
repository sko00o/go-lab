package main

func main() {
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
				println("innner case 0")
				// fallthrough
			case 1:
				println("innner case 1")
				continue Loop
			case 2:
				println("innner case 2")
				break Switch
			}
			println("end for loop")
		}
		println("end case 1")
	}

	println("after switch")
}

/*
before switch
begin case 1
begin for loop
innner case 0
end for loop
begin for loop
innner case 1
end case 1
after switch
*/
/*
before switch
begin case 1
begin for loop
innner case 0
end for loop
begin for loop
innner case 1
begin for loop
innner case 2
after switch
*/
