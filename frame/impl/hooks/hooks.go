package hooks

import (
	"fmt"
	"github.com/opesun/jsonp"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"reflect"
	"strings"
)

// Used to call subscribed hooks.
type Hooks struct {
	hooks     map[string]interface{}
	newModule func(string) iface.Module
	initer    func(iface.Instance) error
	cache     map[string]iface.Instance // Module instance cache.
}

type Hook struct {
	*Hooks
	Hookname string
}

func New(hooks map[string]interface{}, newModule func(string) iface.Module) *Hooks {
	if hooks == nil {
		hooks = map[string]interface{}{}
	}
	return &Hooks{
		hooks,
		newModule,
		nil,
		map[string]iface.Instance{},
	}
}

func (e *Hooks) Initer(initer func(iface.Instance) error) {
	e.initer = initer
}

func (e *Hooks) Select(Hookname string) iface.Hook {
	return &Hook{
		e,
		Hookname,
	}
}

type subscriber struct {
	modName    string
	methodName string
}

func (s *subscriber) Name() string {
	return s.modName
}

func (s *subscriber) Method() string {
	return s.methodName
}

// Return all hooks modules subscribed to a path.
func (e *Hook) subscribers() []*subscriber {
	ret := []*subscriber{}
	subscribed, ok := jsonp.GetS(e.hooks, e.Hookname)
	if !ok {
		return ret
	}
	for _, v := range subscribed {
		inf := &subscriber{}
		switch t := v.(type) {
		case string:
			inf.modName = t
		case []interface{}:
			if len(t) != 2 {
				panic("Misconfigured hook.")
			}
			inf.modName = t[0].(string)
			inf.methodName = t[1].(string)
		}
		ret = append(ret, inf)
	}
	return ret
}

func (e *Hook) Subscribers() []iface.Subscriber {
	s := e.subscribers()
	ret := []iface.Subscriber{}
	for _, v := range s {
		ret = append(ret, v)
	}
	return ret
}

// This is an iface.Module, which wraps the github.com/opesun/nocrud/frame/mod implementation and implements instance caching.
type InstanceCacher struct {
	iface.Module
	cache  map[string]iface.Instance
	initer func(iface.Instance) error
	name   string
}

func (m InstanceCacher) Instance() iface.Instance {
	var ins iface.Instance
	ins, has := m.cache[m.name]
	if !has {
		if !m.Exists() {
			panic(fmt.Sprintf("Module %v does not exist.", m.name))
		}
		insta := m.Module.Instance()
		m.initer(insta)
		m.cache[m.name] = insta
		return insta
	}
	return ins
}

func (e *Hooks) Module(modname string) iface.Module {
	return &InstanceCacher{
		e.newModule(modname),
		e.cache,
		e.initer,
		modname,
	}
}

func (e *Hook) HasSubscribers() bool {
	subscribed := e.subscribers()
	return len(subscribed) > 0
}

func (e *Hook) SubscriberCount() int {
	subscribed := e.subscribers()
	return len(subscribed)
}

// Fire calls hooks subscribed to Hookname, but does not case about their return values.
func (e *Hook) Fire(params ...interface{}) {
	e.iterate(e.Hookname, nil, params...)
}

// Calls all hooks subscribed to Hookname, with params, feeding the output of every hook into stopfunc.
// Stopfunc's argument signature must match the signatures of return values of the called hooks.
// Stopfunc must return a boolean value. A boolean value of true stops the iteration.
// Iterate allows to mimic the semantics of calling all hooks one by one.
func (e *Hook) Iterate(stopfunc interface{}, params ...interface{}) {
	e.iterate(e.Hookname, stopfunc, params...)
}

func validateStopFunc(s reflect.Type) error {
	if s.Kind() != reflect.Func {
		return fmt.Errorf("Stopfunc is not a function.")
	}
	if s.NumOut() != 1 {
		return fmt.Errorf("Stopfunc must have one return value.")
	}
	if s.Out(0) != reflect.TypeOf(false) {
		return fmt.Errorf("Stopfunc must have a boolean return value.")
	}
	return nil
}

func (e *Hook) instance(modname string) iface.Instance {
	ins, exists := e.cache[modname]
	if exists {
		return ins
	}
	mo := e.newModule(modname)
	if !mo.Exists() {
		panic(fmt.Sprintf("Module %v does not exist.", modname))
	}
	insta := mo.Instance()
	if e.initer != nil {
		e.initer(insta)
	}
	e.cache[modname] = insta
	return insta
}

func (e *Hook) iterate(hookName string, stopfunc interface{}, params ...interface{}) {
	subscribed := e.subscribers()
	if len(subscribed) == 0 {
		return
	}
	var stopfunc_numin int
	if stopfunc != nil {
		s := reflect.TypeOf(stopfunc)
		err := validateStopFunc(s)
		if err != nil {
			panic(err)
		}
		stopfunc_numin = s.NumIn()
	}
	nameized := hooknameize(hookName)
	for _, hinf := range subscribed {
		if hinf.methodName == "" {
			hinf.methodName = nameized
		}
		ins := e.instance(hinf.modName)
		if !ins.HasMethod(hinf.methodName) {
			panic(fmt.Sprintf("Module %v has no method named %v", hinf.modName, hinf.methodName))
		}
		var ret_rec interface{}
		hook_outp := []reflect.Value{}
		if stopfunc != nil {
			ret_rec = func(i ...interface{}) {
				for i, v := range i {
					if v == nil {
						hook_outp = append(hook_outp, reflect.Zero(reflect.TypeOf(stopfunc).In(i)))
					} else {
						hook_outp = append(hook_outp, reflect.ValueOf(v))
					}
				}
			}
		}
		err := ins.Method(hinf.methodName).Call(ret_rec, params...)
		if err != nil {
			panic(err)
		}
		if stopfunc != nil {
			if stopfunc_numin != len(hook_outp) {
				panic(fmt.Sprintf("The number of return values of Hook %v of %v differs from the number of arguments of stopfunc.", hinf.methodName, hinf.modName)) // This sentence...
			}
			stopf := reflect.ValueOf(stopfunc)
			stopf_ret := stopf.Call(hook_outp)
			if stopf_ret[0].Interface().(bool) == true {
				break
			}
		}
	}
}

// Creates a hookname from access path.
// "content.insert" => "ContentInsert"
func hooknameize(s string) string {
	s = strings.Replace(s, ".", " ", -1)
	s = strings.Title(s)
	return strings.Replace(s, " ", "", -1)
}
