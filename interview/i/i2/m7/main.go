package main

import (
	"fmt"
	"time"
)

type Project struct{}

func (p *Project) deferError() {
	if err := recover(); err != nil {
		fmt.Println("recover: ", err)
	}
}

func (p *Project) exec(msgchan chan interface{}) {
	defer p.deferError() // fix add recover here

	for msg := range msgchan {
		m := msg.(int) // 断言要用 ok-idioms
		fmt.Println("msg: ", m)
	}
}

func (p *Project) run(msgchan chan interface{}) {
	for {
		// defer p.deferError() // 异常处理应在引发异常的函数中处理
		go p.exec(msgchan)
		time.Sleep(time.Second * 2)
	}
}

func (p *Project) Main() {
	a := make(chan interface{}, 100)
	go p.run(a)
	go func() {
		for {
			a <- "1"
			time.Sleep(time.Second)
		}
	}()
	time.Sleep(time.Second * 100)
}

func main() {
	p := new(Project)
	p.Main()
}
