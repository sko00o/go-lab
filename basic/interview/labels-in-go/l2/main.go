package main

func main() {

Loop:
	for i := 0; i < 3; i++ {
		println("outer for begin")
		for j := 0; j < 3; j++ {
			println("i", i, "j", j)
			switch j {
			case 0:
				println("continue")
				continue
			case 1:
				println("continue Loop")
				continue Loop
			}
			println("inner for end")
		}
		println("outer for end")
	}

}

/*
outer for begin
i 0 j 0
continue
i 0 j 1
continue Loop
outer for begin
i 1 j 0
continue
i 1 j 1
continue Loop
outer for begin
i 2 j 0
continue
i 2 j 1
continue Loop
*/
