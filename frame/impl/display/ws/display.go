package display

import (
	iface "github.com/opesun/nocrud/frame/interfaces"
	"io"
)

type Display struct {
	w    io.Writer
	root string
	typ  string
}

func New(w io.Writer, root string) *Display {
	return &Display{
		w:    w,
		root: root,
	}
}

func (d *Display) Writer() iface.Writer {
	return d.w
}

func (d *Display) Write(dat []byte) error {
	am, err := d.w.Write(dat)
	if am != len(dat) {
		panic("No good.")
	}
	return err
}

func (d *Display) Type(s string) iface.Display {
	return d
}
