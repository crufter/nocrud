package file_return_example

import(
	iface "github.com/opesun/nocrud/frame/interfaces"
)

type C struct {
	ctx		iface.Context
}

func (c *C) Init(ctx iface.Context) {
	c.ctx = ctx
}

func (c *C) FileReturnExample() (iface.File, error) {
	dir, err := c.ctx.FileSys().SelectPlace("modules")
	if err != nil {
		return nil, err
	}
	return dir.File("file_return_example", "tpl", "example.txt"), nil
}