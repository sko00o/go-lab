package main

import (
	"github.com/brocaar/lorawan"
	log "github.com/sirupsen/logrus"
)

type Context struct {
	logger *log.Entry
}

type ops struct {
	a      [3]byte
	DevEUI lorawan.EUI64
}

func main() {
	sub(ops{
		a:      [3]byte{1, 2, 3},
		DevEUI: lorawan.EUI64{1, 2, 3, 4, 5, 6, 7, 8},
	})
}

func sub(o ops) {
	ctx := &Context{
		logger: log.WithFields(log.Fields{
			"a":       o.a,
			"dev_eui": o.DevEUI,
		}),
	}

	for i := 0; i < 3; i++ {
		show(ctx, i)
	}
}

func show(ctx *Context, msg interface{}) {
	logger := ctx.logger
	logger.Info(msg)
}
