package context

import(
	iface "github.com/opesun/nocrud/frame/interfaces"
)

type Context struct {
	conducting		iface.Conducting
	fileSys			iface.FileSys
	user			iface.User
	client	 		iface.Client
	db				iface.Db
	channels		iface.Channels
	viewContext		iface.ViewContext
	nonPortable		iface.NonPortable
	display			iface.Display
	options			iface.Options
}

func New(	cond 	iface.Conducting,
			f 		iface.FileSys,
			u 		iface.User,
			c 		iface.Client,
			db 		iface.Db,
			ch 		iface.Channels,
			v 		iface.ViewContext,
			n 		iface.NonPortable,
			d		iface.Display,
			o		iface.Options,
) iface.Context {
	return &Context{
		cond, f, u, c, db, ch, v, n, d, o,
	}
}

func (c *Context) Conducting() iface.Conducting {
	return c.conducting
}

func (c *Context) FileSys() iface.FileSys {
	return c.fileSys
}

func (c *Context) User() iface.User {
	return c.user
}

func (c *Context) Client() iface.Client {
	return c.client
}

func (c *Context) Db() iface.Db {
	return c.db
}

func (c *Context) Channels() iface.Channels {
	return c.channels
}

func (c *Context) ViewContext() iface.ViewContext {
	return c.viewContext
}

func (c *Context) NonPortable() iface.NonPortable {
	return c.nonPortable
}

func (c *Context) Display() iface.Display {
	return c.display
}

func (c *Context) Options() iface.Options {
	return c.options
}