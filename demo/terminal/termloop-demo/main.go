package main

import (
	tl "github.com/JoelOtter/termloop"
)

func main() {
	game := tl.NewGame()
	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorGreen,
		Fg: tl.ColorBlack,
		Ch: 'v',
	})
	level.AddEntity(tl.NewRectangle(10, 10, 50, 20, tl.ColorBlue))
	player := Player{tl.NewEntity(1,1,1,1)}
	player.SetCell(0,0,&tl.Cell{Fg: tl.ColorRed, Ch: '옷'})
	level.AddEntity(&player)
	game.Screen().SetLevel(level)
	game.Start()
}

type Player struct {
	*tl.Entity
}

func (p *Player) Tick(ev tl.Event) {
	if ev.Type == tl.EventKey {
		x, y := p.Position()
		switch ev.Key {
		case tl.KeyArrowRight:
			p.SetPosition(x+1, y)
		case tl.KeyArrowLeft:
			p.SetPosition(x-1, y)
		case tl.KeyArrowUp:
			p.SetPosition(x, y-1)
		case tl.KeyArrowDown:
			p.SetPosition(x, y+1)
		}
	}
}
