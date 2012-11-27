package filesys

import (
	"fmt"
	"github.com/opesun/nocrud/frame/impl/directory"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/frame/misc/scut"
	"path/filepath"
	"strings"
)

type FileSys struct {
	root string
	host string
	opt  map[string]interface{}
	t    iface.Temporaries
}

func New(root string, host string, opt map[string]interface{}, t iface.Temporaries) *FileSys {
	return &FileSys{
		root,
		host,
		opt,
		t,
	}
}

func dirNameize(s string) string {
	return strings.Replace(s, ":", "-", -1)
}

func (f *FileSys) SelectPlace(s string) (iface.Directory, error) {
	var path string
	switch s {
	case "template":
		path = scut.GetTPath(f.opt, dirNameize(f.host))
	case "modules":
		path = "modules"
	case "uploads":
		path = filepath.Join("uploads", dirNameize(f.host))
	case "root":
		path = ""
	default:
		return nil, fmt.Errorf("Can't find.")
	}
	path = filepath.Join(f.root, path)
	return directory.New(path), nil
}

func (f *FileSys) Temporaries() iface.Temporaries {
	return f.t
}
