package directory_test

import(
	"github.com/opesun/nocrud/frame/impl/directory"
	"testing"
	"os"
	"path/filepath"
)

func exists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil && fi.IsDir() {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func TestCreate(t *testing.T) {
	path := "c:/gowork/renExample"
	os.Remove(path)
	d := directory.New(path)
	err := d.Create()
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
	path := "c:/gowork/renExample"
	os.Remove(path)
	d := directory.New(path)
	err := d.Create()
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
	newname := "laliboy"
	err = d.Rename(newname)
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