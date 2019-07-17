package benchmark

import "testing"

func benchmarkFib(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		Fib(i)
	}
}

// one time loop
// go test -run=^$ -bench "^(BenchmarkFib40)$" -v -benchtime=1ms

func BenchmarkFib1(b *testing.B) {
	benchmarkFib(1, b)
}

func BenchmarkFib5(b *testing.B) {
	benchmarkFib(5, b)
}

func BenchmarkFib10(b *testing.B) {
	benchmarkFib(10, b)
}

func BenchmarkFib40(b *testing.B) {
	benchmarkFib(40, b)
}

var result int

// any benchmark should be careful to avoid compiler optimisations
// eliminating the function under test and artificially lowering
// the run time of the benchmark.
func BenchmarkFibComplete(b *testing.B) {
	var r int
	for n := 0; n < b.N; n++ {
		r = Fib(10)
	}
	result = r
}

func TestFib(t *testing.T) {
	var fibTests = []struct {
		n        int // input
		expected int // expected result
	}{
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 5},
		{6, 8},
		{7, 13},
	}
	for _, tt := range fibTests {
		actual := Fib(tt.n)
		if actual != tt.expected {
			t.Errorf("Fib(%d): expected %d, actual %d", tt.n, tt.expected, actual)
		}
	}
}
