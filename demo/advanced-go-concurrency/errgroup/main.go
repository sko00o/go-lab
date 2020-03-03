package main

import (
	"fmt"

	"golang.org/x/sync/errgroup"
)

// 返回其中任意一个 goroutine 产生的错误

func main() {
	var group errgroup.Group

	for i := 0; i < 12; i++ {
		groupProcess(&group, i)
	}

	if err := group.Wait(); err != nil {
		fmt.Println("Err:", err)
		return
	}

	fmt.Println("no err")
}

func groupProcess(group *errgroup.Group, idx int) {
	group.Go(func() error {
		return noMoreThanNine(idx)
	})
}

func noMoreThanNine(number int) error {
	if number > 9 {
		return fmt.Errorf("%d error", number)
	}

	return nil
}
