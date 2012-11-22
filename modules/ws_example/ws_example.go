package ws_example

import (
	"fmt"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"time"
)

type C struct {
	d iface.Display
}

type Ctx interface {
	Display() iface.Display
}

func (c *C) Init(ctx Ctx) {
	c.d = ctx.Display()
}

func (c *C) WsHello() {
	for {
		err := c.d.Write([]byte("Hello there."))
		if err != nil {
			break
		}
		time.Sleep(time.Second)
	}
	fmt.Println("Quitting.")
}
