package cp

import (
	"fmt"
	"os"
	"runtime/trace"
	"testing"
)

func BenchmarkMemoryStack(b *testing.B) {
	var s S

	f, err := os.Create("stack.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		s = byCopy()
	}

	trace.Stop()

	b.StopTimer()

	_ = fmt.Sprintf("%v", s.a)
}

func BenchmarkMemoryHeap(b *testing.B) {
	var s *S

	f, err := os.Create("heap.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		s = byPointer()
	}

	trace.Stop()

	b.StopTimer()

	_ = fmt.Sprintf("%v", s.a)
}

func BenchmarkMemoryStack2(b *testing.B) {
	var s, s1 S
	s = byCopy()
	s1 = byCopy()
	for i := 0; i < b.N; i++ {
		for i := 0; i < 1000000; i++ {
			s.stack(s1)
		}
	}
}

func BenchmarkMemoryHeap2(b *testing.B) {
	var s, s1 *S
	s = byPointer()
	s1 = byPointer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < 1000000; i++ {
			s.heap(s1)
		}
	}
}
