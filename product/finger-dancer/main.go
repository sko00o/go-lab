package main

import (
	"strings"

	"github.com/nsf/termbox-go"
)

var testStr = "this is a test string. this is a test string. this is a test string. this is a test string. this is a test string. this is a test string. this is a test string. this is a test string. this is a test string. this is a test string. this is a test string. this is a test string. this is a test string. this is a test string. this is a test string. this is a test string. this is a test string. this is a test string."

func main() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	keys := initStr(testStr)
	keys.draw()

	var cIdx int
loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			// hit Ctrl+Q to quit
			if ev.Key == termbox.KeyCtrlQ {
				break loop
			}

			if cIdx < len(keys) {
				if judge(&ev, &keys[cIdx]) {
					keys[cIdx].s = Right
				} else {
					keys[cIdx].s = Wrong
				}

				cIdx++
				if cIdx < len(keys) {
					keys[cIdx].s = Current
				}
			}

			keys.draw()

		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

func initStr(str string) (out keys) {
	var xIdx, yIdx int
	for idx, r := range str {

		if wEnd := xIdx + strings.Index(str[idx:], " "); wEnd >= 80 {
			yIdx++
			xIdx = 0
		}
		k := key{
			x:  xIdx,
			y:  yIdx,
			ch: r,
		}

		if idx == 0 {
			k.s = Current
		}
		out = append(out, k)
		xIdx++
	}
	return
}

func judge(ev *termbox.Event, k *key) bool {
	if ev.Ch == k.ch {
		return true
	} else if ev.Ch == 0 &&
		ev.Key == termbox.KeySpace &&
		k.ch == rune(' ') {
		return true
	}
	return false
}

type key struct {
	x, y int
	ch   rune

	s statInt
}

type keys []key

func (k keys) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for _, e := range k {
		s := e.s.Stat()
		termbox.SetCell(e.x, e.y, e.ch, s.fg, s.bg)
	}
	termbox.Flush()
}

type statInt int

const (
	Unknown statInt = iota
	Wrong
	Right
	Current
)

type stat struct {
	fg termbox.Attribute
	bg termbox.Attribute
}

func (s statInt) Stat() *stat {
	switch s {
	case Unknown:
		return &stat{fg: termbox.ColorWhite, bg: termbox.ColorBlack}
	case Wrong:
		return &stat{fg: termbox.ColorYellow, bg: termbox.ColorRed}
	case Right:
		return &stat{fg: termbox.ColorGreen, bg: termbox.ColorBlack}
	case Current:
		return &stat{fg: termbox.ColorWhite, bg: termbox.ColorBlue}
	}
	return &stat{fg: termbox.ColorDefault, bg: termbox.ColorDefault}
}
