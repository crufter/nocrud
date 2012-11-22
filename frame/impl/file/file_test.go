package file_test

import (
	"github.com/opesun/nocrud/frame/impl/file"
	"os"
	"path/filepath"
	"testing"
)

func exists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil && !fi.IsDir() {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func TestCreate(t *testing.T) {
	path := "c:/gowork/example.txt"
	os.Remove(path)
	f := file.New(path)
	err := f.Create()
	if err != nil {
		t.Fatal(err)
	}
	ex, err := exists(path)
	if err != nil {
		t.Fatal(err)
	}
	if !ex {
		t.Fatal("Not created.")
	}
}

func TestRename(t *testing.T) {
	path := "c:/gowork/example.txt"
	os.Remove(path)
	f := file.New(path)
	err := f.Create()
	if err != nil {
		t.Fatal(err)
	}
	ex, err := exists(path)
	if err != nil {
		t.Fatal(err)
	}
	if !ex {
		t.Fatal("Not created.")
	}
	newname := "laliboy.txt"
	err = f.Rename(newname)
	if err != nil {
		t.Fatal(err)
	}
	newpath := filepath.Join(filepath.Dir(path), newname)
	ex, err = exists(newpath)
	if err != nil {
		t.Fatal(err)
	}
	if !ex {
		t.Fatal("Rename didn't happen.")
	}
}
