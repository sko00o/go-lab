package main

type humanImpl struct {
	base
}

func (b *humanImpl) Walk() int {
	return 123
}

func newHuman() human {
	return &humanImpl{
		base: newBase(),
	}
}
