package nesteddata

import (
	"github.com/opesun/jsonp"
	"github.com/opesun/numcon"
	"strings"
)

type NestedData struct {
	d interface{}
}

func New(n interface{}) *NestedData {
	return &NestedData{
		n,
	}
}

func (n *NestedData) Get(s ...string) (interface{}, bool) {
	p := strings.Join(s, ".")
	return jsonp.Get(n.d, p)
}

func (n *NestedData) GetI(s ...string) (int64, bool) {
	r, ok := n.Get(s...)
	if !ok {
		return 0, false
	}
	val, err := numcon.Int64(r)
	if err != nil {
		return 0, false
	}
	return val, true
}

func (n *NestedData) GetF(s ...string) (float64, bool) {
	r, ok := n.Get(s...)
	if !ok {
		return 0, false
	}
	val, err := numcon.Float64(r)
	if err != nil {
		return 0, false
	}
	return val, true
}

func (n *NestedData) GetS(s ...string) ([]interface{}, bool) {
	r, ok := n.Get(s...)
	val, ok := r.([]interface{})
	return val, ok
}

func (n *NestedData) GetStr(s ...string) (string, bool) {
	r, ok := n.Get(s...)
	val, ok := r.(string)
	return val, ok
}

func (n *NestedData) GetM(s ...string) (map[string]interface{}, bool) {
	r, ok := n.Get(s...)
	val, ok := r.(map[string]interface{})
	return val, ok
}

func (n *NestedData) Exists(s ...string) bool {
	_, ok := n.Get(s...)
	return ok
}

func (n *NestedData) All() interface{} {
	return n.d
}
