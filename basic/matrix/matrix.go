package matrix

type Matrix [][]int64

func New(size int) Matrix {
	res := make([][]int64, size)
	for i := 0; i < size; i++ {
		res[i] = make([]int64, size)
	}
	return res
}

func (m1 Matrix) canAdd(m2 Matrix) int {
	size := len(m1)
	if size != len(m2) {
		panic("require same size")
	}
	return size
}

func (m1 Matrix) AddRowToRow(m2 Matrix) {
	size := m1.canAdd(m2)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			m1[i][j] = m1[i][j] + m2[i][j]
		}
	}
}

func (m1 Matrix) AddColumnToRow(m2 Matrix) {
	size := m1.canAdd(m2)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			m1[i][j] = m1[i][j] + m2[j][i]
		}
	}
}

// loop nest optimization
func (m1 Matrix) AddColumnToRawBlock(m2 Matrix, blockSize int) {
	size := m1.canAdd(m2)
	for i := 0; i < size; i += blockSize {
		for j := 0; j < size; j += blockSize {
			for ii := i; ii < i+blockSize && ii < size; ii++ {
				for jj := j; jj < j+blockSize && jj < size; jj++ {
					m1[ii][jj] = m1[ii][jj] + m2[jj][ii]
				}
			}
		}
	}
}
