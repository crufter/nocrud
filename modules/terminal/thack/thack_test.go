package thack_test

import(
	"github.com/opesun/nocrud/modules/terminal/thack"
	"testing"
	"fmt"
)

func TestBasic(t *testing.T) {
	f := map[string]interface{}{
		"testFunc": func() string {
			return "Hello."
		},
	}
	script := "testFunc"
	s, err := thack.New().Funcs(f).Execute(script)
	if err != nil {
		t.Fatal(err)
	}
	if s != "Hello." {
		t.Fatal(s)
	}
}

func TestComments(t *testing.T) {
	f := map[string]interface{}{
		"testFunc": func() string {
			return "Hello."
		},
	}
	script := "testFunc\n#Comment stuff\ntestFunc"
	s, err := thack.New().Funcs(f).Execute(script)
	if err != nil {
		t.Fatal(err)
	}
	if s != "Hello.\nHello." {
		t.Fatal(s, len(s))
	}
}

func TestError(t *testing.T) {
	f := map[string]interface{}{
		"testE": func() error {
			return fmt.Errorf("Test error.")
		},
		"t1": func() string {
			return "Hello."
		},
	}
	script := "$x := testE\n$x\nt1"
	s, err := thack.New().Funcs(f).Execute(script)
	if err == nil {
		t.Fatal("Error was not returned to Execute.")
	}
	if s != "" {
		t.Fatal(s, len(s))
	}
}