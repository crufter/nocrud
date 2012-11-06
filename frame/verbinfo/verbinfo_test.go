package verbinfo_test

import(
	"testing"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/frame/verbinfo"
	"github.com/opesun/nocrud/frame/mod"
)

type mockObject struct{}

func (m *mockObject) MethodA(a iface.Filter, data map[string]interface{}) {
	
}

func (m *mockObject) MethodB(a, b iface.Filter) {
	
}

func TestA(t *testing.T) {
	methA := mod.ToInstance(&mockObject{}).Method("MethodA")
	an := verbinfo.NewAnalyzer(methA)
	if an.ArgCount() != 2 {
		t.Fatal()
	}
	if an.FilterCount() != 1 {
		t.Fatal()
	}
	if !an.NeedsData() {
		t.Fatal()
	}
}

func TestB(t *testing.T) {
	methB := mod.ToInstance(&mockObject{}).Method("MethodB")
	an := verbinfo.NewAnalyzer(methB)
	if an.ArgCount() != 2 {
		t.Fatal()
	}
	if an.FilterCount() != 2 {
		t.Fatal()
	}
	if an.NeedsData() {
		t.Fatal()
	}
}