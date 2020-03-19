package buffer

import (
	"math/rand"
	"testing"
)

// Manual Loop Faster Than `range` Operator

var s1, s2 []uint16

const (
	listScale = 500
)

func init() {
	s1 = make([]uint16, listScale*2)
	s2 = make([]uint16, listScale)

	val := uint16(0)
	for i := 0; i < len(s1); i++ {
		val += uint16(rand.Intn(3) + 1)
		s1[i] = val
	}

	val = uint16(0)
	for i := 0; i < len(s2); i++ {
		val += uint16(rand.Intn(6) + 1)
		s2[i] = val
	}
}

func BenchmarkHandRange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HandRange(s1, s2)
		HandRange(s2, s1)
	}
}

func BenchmarkHandRange2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HandRange(s2, s1)
		HandRange(s1, s2)
	}
}

func BenchmarkIterRange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IterRange(s1, s2)
		IterRange(s2, s1)
	}
}

func BenchmarkIterRange2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IterRange(s2, s1)
		IterRange(s1, s2)
	}
}
