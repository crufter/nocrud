package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type File struct {
	path string
}

func New(path string) *File {
	return &File{
		path,
	}
}

func (f *File) Create() error {
	err := os.MkdirAll(filepath.Dir(f.path), os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f.path, []byte{}, os.ModePerm)
}

func (f *File) Exists() (bool, error) {
	fi, err := os.Stat(f.path)
	if err == nil && !fi.IsDir() {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (f *File) Read() ([]byte, error) {
	return ioutil.ReadFile(f.path)
}

func (f *File) Write(b []byte) error {
	return ioutil.WriteFile(f.path, b, os.ModePerm)
}

func (f *File) Remove() error {
	return os.Remove(f.path)
}

// Copied from github.com/opesun/frame/impl/directory ... what could I do...
func (f *File) Rename(newName string) error {
	newPath := filepath.Join(filepath.Dir(f.path), newName)
	err := os.Rename(f.path, newPath)
	if err != nil {
		return err
	}
	f.path = newPath
	return nil
}

func (f *File) Name() string {
	return filepath.Base(f.path)
}
