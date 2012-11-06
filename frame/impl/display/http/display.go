package display

import(
	iface "github.com/opesun/nocrud/frame/interfaces"
	"net/http"
)

type Display struct {
	w		http.ResponseWriter
	root	string
	typ		string
}

func New(w http.ResponseWriter, root string) *Display {
	return &Display{
		w:		w,
		root:	root,
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
	switch s {
	case "html":
		d.w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case "json":
		d.w.Header().Set("Content-Type", "application/json; charset=utf-8")
	default:
		panic("Unkown content type " + s + ".")
	}
	return d
}