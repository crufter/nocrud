package directory

import(
	"os"
	"io/ioutil"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/frame/impl/file"
	"path/filepath"
)

type Directory struct {
	path	string
}

func New(path string) *Directory {
	return &Directory{
		path,
	}
}

func (d *Directory) Exists() (bool, error) {
	fi, err := os.Stat(d.path)
	if err == nil && fi.IsDir() {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (d *Directory) Remove() error {
	return os.Remove(d.path)
}

type FileInfo struct {
	name	string
	path	string
	isDir	bool
}

func (f *FileInfo) File() iface.File {
	if f.isDir {
		panic("Not a file")
	}
	return file.New(f.path)
}

func (f *FileInfo) IsDir() bool {
	return f.isDir
}

func (f *FileInfo) Name() string {
	return f.name
}

func (f *FileInfo) Directory() iface.Directory {
	if !f.isDir {
		panic("Not a directory.")
	}
	return &Directory{
		f.path,
	}
}

func (d *Directory) List() ([]iface.FileInfo, error) {
	finfos, err := ioutil.ReadDir(d.path)
	if err != nil {
		return nil, err
	}
	ret := []iface.FileInfo{}
	for _, v := range finfos {
		ret = append(ret, &FileInfo{
			name:	v.Name(),
			path:	filepath.Join(d.path, v.Name()),
			isDir:	v.IsDir(),
		})
	}
	return ret, nil
}

func (d *Directory) Create() error {
	return os.MkdirAll(d.path, os.ModePerm)
}

func (d *Directory) Rename(newName string) error {
	newPath := filepath.Join(filepath.Dir(d.path), newName)
	err := os.Rename(d.path, newPath)
	if err != nil {
		return err
	}
	d.path = newPath
	return nil
}

func (d *Directory) Directory(s ...string) iface.Directory {
	arg := []string{d.path}
	arg = append(arg, s...)
	return New(filepath.Join(arg...))
}

func (d *Directory) File(s string) iface.File {
	return file.New(filepath.Join(d.path, s))
}