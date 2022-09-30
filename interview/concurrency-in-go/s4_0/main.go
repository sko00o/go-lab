package main

import (
	"fmt"
)

var tasks = []func(){
	ConfinementAdHoc,
	ConfinementLexical,
}

func main() {
	for _, t := range tasks {
		t()
	}
}

func ConfinementAdHoc() {
	data := []int{1, 2, 3, 4}
	loopData := func(handleData chan<- int) {
		defer close(handleData) // may forget to close the channel
		for i := range data {
			handleData <- data[i]
		}
	}

	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}
}

// Recommended
func ConfinementLexical() {
	chanOwner := func() <-chan int {
		results := make(chan int, 5)
		go func() {
			defer close(results)
			for i := 0; i <= 5; i++ {
				results <- i
			}
		}()
		return results
	}

	consumer := func(results <-chan int) {
		for result := range results {
			fmt.Println(result)
		}
		fmt.Println("done")
	}

	res := chanOwner()
	consumer(res)
}
