package conducting

import (
	iface "github.com/opesun/nocrud/frame/interfaces"
)

type Conducting struct {
	hooks  iface.Hooks
	events iface.Events
}

func New(hooks iface.Hooks, events iface.Events) *Conducting {
	return &Conducting{
		hooks,
		events,
	}
}

func (c *Conducting) Hooks() iface.Hooks {
	return c.hooks
}

func (c *Conducting) Events() iface.Events {
	return c.events
}
