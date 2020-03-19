package tail_recursion

import (
	"testing"
)

func BenchmarkFactorial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fact(20)
	}
}

func BenchmarkFactorialTail(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tailFact(20)
	}
}
