package file

import(
	"os"
	"io/ioutil"
	"fmt"
	"strings"
)

type File struct {
	path	string
}

func New(path string) *File {
	return &File{
		path,
	}
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

func (f *File) Rename(s string) error {
	return fmt.Errorf("!")
}

func (f *File) Name() string {
	spl := strings.Split(f.path, "/")
	return spl[len(spl)-1]
}
