package filesys

import(
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/frame/impl/directory"
	"path/filepath"
	"github.com/opesun/nocrud/frame/misc/scut"
	"fmt"
)

type FileSys struct {
	root	string
	host	string
	opt 	map[string]interface{}
	t		iface.Temporaries
}

func New(root string, host string, opt map[string]interface{}, t iface.Temporaries) *FileSys {
	return &FileSys{
		root,
		host,
		opt,
		t,
	}
}

func (f *FileSys) SelectPlace(s string) (iface.Directory, error) {
	var path string
	switch s {
	case "template":
		path = scut.GetTPath(f.opt, f.host)
	case "modules":
		path = "modules"
	case "uploads":
		path = "uploads"
	default:
		return nil, fmt.Errorf("Can't find.")
	}
	path = filepath.Join(f.root, path)
	return directory.New(path), nil
}

func (f *FileSys) Temporaries() iface.Temporaries {
	return f.t
}