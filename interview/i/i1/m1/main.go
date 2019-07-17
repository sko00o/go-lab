package main

func deferCall() {
	defer print(1)
	defer print(2)
	defer print(3)

	panic("c")
}

func main() {
	deferCall()
}
