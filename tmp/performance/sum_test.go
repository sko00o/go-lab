package sum

import (
	"testing"
)

func BenchmarkSum(b *testing.B) { // 37s
	s := new(int)
	*s = 0
	sum(s)
}

func BenchmarkSum1(b *testing.B) { // 10s
	s := new(int)
	*s = 0
	sum1(s)
}

func BenchmarkSum2(b *testing.B) { // 10s
	b.ResetTimer()
	s := new(int)
	*s = 0
	sum2(s)
	b.StopTimer()
}
