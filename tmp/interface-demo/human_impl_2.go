package main

type humanImpl2 struct {
	base
}

func (b *humanImpl2) Walk() int {
	return 456
}

func (b *humanImpl2) Talk() string {
	return "human"
}

func newHuman21() human {
	return &humanImpl2{
		base: newBase(),
	}
}

func newHuman22() human {
	return &humanImpl2{}
}
