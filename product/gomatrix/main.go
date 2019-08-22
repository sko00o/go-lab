package main

import (
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

type matrix struct {
	val    int
	isHead bool
}

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	count := 0
	update := 4 * time.Millisecond

	cols, rows := termbox.Size()

	min, max := 33, 123
	rNum := max - min

	// 初始化
	m := make([][]matrix, rows)
	for i := range m {
		m[i] = make([]matrix, cols)
		for j := 0; j < cols; j += 2 {
			m[i][j].val = -1
		}
	}
	updates := make([]int, cols) // 决定当前线是否更新
	spaces := make([]int, cols)  // 每列的留白长度
	length := make([]int, cols)  // 每条线的列长度
	for j := 0; j < cols; j += 2 {
		spaces[j] = rand.Int()%rows + 1
		length[j] = rand.Int()%(rows-3) + 3
		m[1][j].val = ' '
		updates[j] = rand.Int()%3 + 1
	}

	quit := false
	go func() {
		for {
			switch ev := termbox.PollEvent(); ev.Key {
			case termbox.KeyEsc:
				quit = true
				return
			}
		}
	}()

	for {
		if quit {
			return
		}

		count++
		if count > 4 {
			count = 1
		}

		for j := 0; j < cols; j += 2 {

			if count > updates[j] { // 决定是否要更新
				// old-style scrolling
				// 所有字符下移一个单位
				for i := rows - 1; i > 0; i-- {
					m[i][j].val = m[i-1][j].val
				}

				random := rand.Int()%(rNum+8) + min

				if m[1][j].val == 0 {
					m[0][j].val = 1
				} else if m[1][j].val == ' ' || // 第一行的 j 列没字符
					m[1][j].val == -1 {

					if spaces[j] > 0 { // 第 j 列还能填空字符， 就继续填
						m[0][j].val = ' '
						spaces[j]--
					} else { // 没得填了，就选随机字符

						// 随机数决定是下一列的头是否有‘白’头
						if rand.Int()%3 == 1 {
							m[0][j].val = 0
						} else {
							m[0][j].val = rand.Int()%rNum + min // 给一个随机字符
						}
						spaces[j] = rand.Int()%rows + 1 // 第 j 列的留白长度 随机更新， 但至少有一个
					}

				} else if random > max && m[1][j].val != 1 {
					m[0][j].val = ' '
				} else {
					m[0][j].val = rand.Int()%rNum + min
				}
			}

			for i := 0; i < rows; i++ {
				termbox.SetCursor(i, j)

				if m[i][j].val == 0 || m[i][j].isHead {
					if m[i][j].val == 0 {
						draw(i, j, '&')
					} else {
						draw(i, j, rune(m[i][j].val))
					}

				} else {

					if m[i][j].val == 1 {
						draw(i, j, '|')
					} else {
						if m[i][j].val == -1 {
							draw(i, j, ' ')
						} else {
							draw(i, j, rune(m[i][j].val))
						}
					}

				}
			}

		}

		termbox.Flush()
		time.Sleep(update * 10)
	}
}

func draw(x, y int, ch rune) {
	termbox.SetCell(y, x, ch, termbox.ColorGreen, termbox.ColorDefault)
}
