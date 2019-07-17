package main

func main() {
	println("before for")
Loop:
	for i := 0; i < 3; i++ {
		println("for begin")
		switch i {
		case 0:
			println("break")
			break
		case 1:
			println("break Loop")
			break Loop
		}
		println("for end")
	}
	println("after for")
}

/*
before for
for begin
break
for end
for begin
break Loop
after for
*/
