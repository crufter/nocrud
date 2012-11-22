package verbinfo

import (
	iface "github.com/opesun/nocrud/frame/interfaces"
	"reflect"
)

type Analyzer struct {
	m iface.Method
}

func NewAnalyzer(method iface.Method) *Analyzer {
	return &Analyzer{method}
}

func (a *Analyzer) ArgCount() int {
	return len(a.m.InputTypes())
}

func (a *Analyzer) FilterCount() int {
	inptypes := a.m.InputTypes()
	c := 0
	var i *iface.Filter
	ft := reflect.TypeOf(i).Elem()
	for _, v := range inptypes {
		if v.Implements(ft) {
			c++ // lol.
		}
	}
	return c
}

func (a *Analyzer) NeedsData() bool {
	inptypes := a.m.InputTypes()
	if len(inptypes) == 0 {
		return false
	}
	return inptypes[len(inptypes)-1] == reflect.TypeOf(map[string]interface{}{})
}

// Return value analyzer...
func NewRanalyzer(a []interface{}) *Ranalyzer {
	return &Ranalyzer{a}
}

type Ranalyzer struct {
	a []interface{}
}

func (r *Ranalyzer) HadError() bool {
	if len(r.a) == 0 {
		return false
	}
	if val, ok := r.a[len(r.a)-1].(error); ok && val != nil {
		return true
	}
	return false
}

func (r *Ranalyzer) Error() error {
	if len(r.a) == 0 {
		return nil
	}
	if val, ok := r.a[len(r.a)-1].(error); ok && val != nil {
		return val
	}
	return nil
}

func (r *Ranalyzer) NonErrorCount() int {
	if len(r.a) == 0 {
		return 0
	}
	c := 0
	for _, v := range r.a {
		_, is_err := v.(error)
		if is_err {
			c++
		}
	}
	return c
}

func (r *Ranalyzer) NonErrors() []interface{} {
	ret := []interface{}{}
	if len(r.a) == 0 {
		return ret
	}
	for _, v := range r.a {
		_, is_err := v.(error)
		if !is_err {
			ret = append(ret, v)
		}
	}
	return ret
}
