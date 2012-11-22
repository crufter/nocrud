package viewcontext_test

import (
	//iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/frame/impl/viewcontext"
	"testing"
)

func TestBasic(t *testing.T) {
	vctx := viewcontext.New()
	vctx.Publish("str", "x")
	vctx.Publish("int", 12)
	m := vctx.Get()
	if m["str"] != "x" {
		t.Fatal()
	}
	if m["int"] != 12 {
		t.Fatal()
	}
}

// Should not handle dot notation.
func TestNonDot(t *testing.T) {
	vctx := viewcontext.New()
	vctx.Publish("whatever.whatever", "asdasd")
	m := vctx.Get()
	_, ok := m["whatever"]
	if ok {
		t.Fatal()
	}
}

// Existing map under a given key + another map published with same key = two maps merged
func TestMerge(t *testing.T) {
	vctx := viewcontext.New()
	tm := map[string]interface{}{
		"x": 1,
	}
	tm1 := map[string]interface{}{
		"y": 2,
	}
	mname := "mapname"
	vctx.Publish(mname, tm)
	vctx.Publish(mname, tm1)
	m := vctx.Get()
	if len(m) != 1 {
		t.Fatal(m)
	}
	if m[mname].(map[string]interface{})["x"] != 1 || m[mname].(map[string]interface{})["y"] != 2 {
		t.Fatal(m)
	}
}
