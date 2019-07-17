package main

type baseImpl struct {
}

func (b *baseImpl) Talk() string {
	return "base"
}

func newBase() base {
	return &baseImpl{}
}
