package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	q := NewSomethingQueue()
	assert.NotNil(t, q)
}

func TestPushAndPop(t *testing.T) {
	ast := assert.New(t)
	item1 := new(Something)
	item2 := new(Something)
	q := NewSomethingQueue()

	q.Push(item1)
	q.Push(item2)

	ast.Equal(item2, q.Pop())
	ast.Equal(item1, q.Pop())
}
