package matrix

import (
	"strconv"
	"testing"
)

func BenchmarkAddRowToRow(b *testing.B) {
	b.StopTimer()
	const size = 6400
	m1 := New(size)
	m2 := New(size)
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		m1.AddRowToRow(m2)
	}
}

func BenchmarkAddColumnToRow(b *testing.B) {
	b.StopTimer()
	const size = 6400
	m1 := New(size)
	m2 := New(size)
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		m1.AddColumnToRow(m2)
	}
}

func BenchmarkAddColumnToRowBlock(b *testing.B) {
	b.StopTimer()
	const size = 6400
	m1 := New(size)
	m2 := New(size)
	b.StartTimer()

	for _, blockSize := range []int{
		16, 32, 64, 128,
	} {
		b.Run("bloclSize="+strconv.Itoa(blockSize), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				m1.AddColumnToRawBlock(m2, blockSize)
			}
		})
	}
}
