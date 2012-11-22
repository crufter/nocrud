package channels

import (
	iface "github.com/opesun/nocrud/frame/interfaces"
)

type Channels struct {
	m map[string][]interface{}
}

func New() *Channels {
	return &Channels{
		map[string][]interface{}{},
	}
}

type Channel struct {
	m    map[string][]interface{}
	name string
}

func (c *Channels) Select(s string) iface.Channel {
	return &Channel{
		c.m,
		s,
	}
}

func (c *Channel) HasData() bool {
	_, ok := c.m[c.name]
	return ok
}

func (c *Channel) Send(i interface{}) {
	c.m[c.name] = append(c.m[c.name], i)
}

func (c *Channel) Get() []interface{} {
	return c.m[c.name]
}

func (c *Channel) GetFirst() interface{} {
	if len(c.m[c.name]) > 0 {
		return c.m[c.name][0]
	}
	return nil
}

func (c *Channel) GetX(i int) interface{} {
	if len(c.m[c.name]) > i {
		return c.m[c.name][i]
	}
	return nil
}

func (c *Channel) HasX(i int) bool {
	if len(c.m[c.name]) > i {
		return true
	}
	return false
}
