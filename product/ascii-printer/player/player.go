package player

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	video2 "github.com/sko00o/go-lab/product/ascii-printer/video"
)

const (
	ESC         = "\033"      // linux escape character
	ClearScreen = ESC + "[2J" // Clear entire screen

	MoveTmpl = ESC + "[%d;%dH"
)

type player struct {
	video  *video2.Video
	writer io.Writer

	screenWidth  int
	screenHeight int

	cursor  time.Duration
	stopped bool

	Speed time.Duration

	timeBar *timeBar
}

func DefaultPlayer(writer io.Writer, video *video2.Video) *player {
	return &player{
		writer: writer,
		video:  video,

		screenWidth:  67,
		screenHeight: 14,

		Speed: 30,

		timeBar: DefaultTimeBar(video.TotalTime),
	}
}

func (p *player) Play() {
	p.stopped = false
	for _, frame := range p.video.Frames {
		if p.stopped {
			return
		}

		p.cursor += frame.DisplayTime
		p.loadFrame(frame, p.cursor)
		if f, ok := p.writer.(http.Flusher); ok {
			f.Flush()
		}
		time.Sleep(frame.DisplayTime * time.Second / p.Speed)
	}
}

func (p *player) Stop() {
	p.stopped = true
}

func (p *player) loadFrame(frame video2.Frame, pos time.Duration) {
	var screenBuf bytes.Buffer

	screenBuf.Write([]byte(ClearScreen))

	screenBuf.Write(p.moveCursor(1, 1))

	for _, line := range frame.Content {
		screenBuf.Write([]byte(line + "\r\n"))
	}

	p.updateTimeBar(&screenBuf, pos)

	// draw
	fmt.Fprint(p.writer, screenBuf.String())
}

func (p *player) updateTimeBar(screenBuf *bytes.Buffer, pos time.Duration) {
	screenBuf.Write(p.moveCursor(1, p.screenHeight))
	screenBuf.Write([]byte(GetTimeBar(pos) + "\r\n"))
}

func (p *player) moveCursor(x, y int) []byte {
	if x <= 0 || y <= 0 ||
		x > p.screenWidth ||
		y > p.screenHeight {
		log.Printf("warning, coordinates out of range. (%d, %d)\n", x, y)
		return nil
	}

	return []byte(fmt.Sprintf(MoveTmpl, y, x))
}
