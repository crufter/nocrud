package display

// All functions which can be called from templates reside here.

import (
	"github.com/opesun/nocrud/frame/lang"
	"github.com/opesun/nocrud/frame/misc/convert"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/frame/highlev"
	"github.com/opesun/nocrud/frame/misc/scut"
	"github.com/opesun/jsonp"
	"github.com/opesun/paging"
	"github.com/opesun/numcon"
	"html/template"
	"reflect"
	"strings"
	"time"
	"strconv"
	"fmt"
)

func get(dat map[string]interface{}, s ...string) interface{} {
	if len(s) > 0 {
		if len(s[0]) > 0 {
			if string(s[0][0]) == "$" {
				s[0] = s[0][1:]
			}
		}
	}
	access := strings.Join(s, ".")
	val, has := jsonp.Get(dat, access)
	if !has {
		return access
	}
	return val
}

func date(timestamp int64, format ...string) string {
	var form string
	if len(format) == 0 {
		form = "2006.01.02 15:04:05"
	} else {
		form = format[0]
	}
	t := time.Unix(timestamp, 0)
	return t.Format(form)
}

func isMap(a interface{}) bool {
	v := reflect.ValueOf(a)
	switch kind := v.Kind(); kind {
	case reflect.Map:
		return true
	}
	return false
}

func eq(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func html(s string) template.HTML {
	return template.HTML(s)
}

func nonEmpty(a interface{}) bool {
	if a == nil {
		return false
	}
	switch t := a.(type) {
	case string:
		return t != ""
	case bool:
		return t != false
	default:
		return true
	}
	return true
}

// Returns the first argument which is not nil, false or empty string.
// Returns false if none of the arguments matches that criteria.
func fallback(a ...interface{}) interface{} {
	for _, v := range a {
		if nonEmpty(v) {
			return v
		}
	}
	return ""
}

func formatFloat(i interface{}, prec int) string {
	f, err := numcon.Float64(i)
	if err != nil {
		return err.Error()
	}
	return strconv.FormatFloat(f, 'f', prec, 64)
}

// For debugging purposes.
func typeOf(i interface{}) string {
	return fmt.Sprint(reflect.TypeOf(i))
}

func sameKind(a, b interface{}) bool {
	return reflect.ValueOf(a).Kind() == reflect.ValueOf(b).Kind()
}

type Form struct {
	*lang.Form
}

func (f *Form) HiddenFields() [][2]string {
	ret := [][2]string{}
	for i, v := range f.FilterFields {
		if _, yepp := v.([]interface{}); yepp {
			for _, x := range v.([]interface{}) {
				ret = append(ret, [2]string{i, fmt.Sprint(x)})
			}
		} else {
			ret = append(ret, [2]string{i, fmt.Sprint(v)})
		}
	}
	return ret
}

func (f *Form) HiddenString() template.HTML {
	d := f.HiddenFields()
	ret := ""
	for _, v := range d {
		ret = ret+`<input type="hidden" name="`+v[0]+`" value="`+v[1]+`" />`
	}
	return template.HTML(ret)
}

func form(ctx iface.Context, action_name string) *Form {
	nouns := scut.GetNouns(ctx.Options().Document())
	np := ctx.NonPortable()
	hl, err := highlev.New(ctx.Conducting().Hooks(), np.Resource(), nouns, np.Params())
	if err != nil {
		panic(err)
	}
	f := hl.URLE().Form(action_name)
	return &Form{
		f,
	}
}

func _url(ctx iface.Context, action_name string, i ...interface{}) string {
	if len(i)%2 == 1 {
		panic("Must be even.")
	}
	nouns := scut.GetNouns(ctx.Options().Document())
	np := ctx.NonPortable()
	hl, err := highlev.New(ctx.Conducting().Hooks(), np.Resource(), nouns, np.Params())
	if err != nil {
		panic(err)
	}
	f := hl.URLE()
	inp := convert.ListToMap(i...)
	return f.UrlString(action_name, inp)
}

type counter int

func newcounter() *counter {
	v := counter(0)
	return &v
}

func (c *counter) Inc() string {		// Ugly hack, template engine needs a return value.
	*c++
	return ""
}

func (c counter) Eq(i int) bool {
	return int(c) == i
}

func (c counter) EveryX(i int) bool {
	if i == 0 {
		return false
	}
	return int(c)%i==0
}

// Works from Get or GetSingle only.
func getSub(ctx iface.Context, noun string, params ...interface{}) []interface{} {
	nouns, ok := ctx.Options().Document().GetM("nouns")
	if !ok {
		panic("Can't find nouns.")
	}
	np := ctx.NonPortable()
	hl, err := highlev.New(ctx.Conducting().Hooks(), np.Resource(), nouns, np.Params())
	if err != nil {
		panic(err)
	}
	inp := convert.ListToMap(params...)
	subhl, err := hl.Sub(noun, inp)
	if err != nil {
		panic(err)
	}
	deflev_i, _ := ctx.Options().Document().Get("default_level")
	deflev, _ := numcon.Int(deflev_i)
	ret, err := subhl.Run(ctx.Db(), ctx.User(), deflev)
	if err != nil {
		panic(err)
	}
	return ret
}

func getList(ctx iface.Context, noun string, params ...interface{}) []interface{} {
	nouns, ok := ctx.Options().Document().GetM("nouns")
	if !ok {
		panic("Can't find nouns.")
	}
	inp := convert.ListToMap(params...)
	hl, err := highlev.New(ctx.Conducting().Hooks(), "/" + noun, nouns, inp)
	if err != nil {
		panic(err)
	}
	deflev_i, _ := ctx.Options().Document().Get("default_level")
	deflev, _ := numcon.Int(deflev_i)
	ret, err := hl.Run(ctx.Db(), ctx.User(), deflev)
	if err != nil {
		panic(err)
	}
	return ret
}

type pagr struct {
	HasPrev		bool
	Prev		int
	HasNext		bool
	Next		int
	Elems		[]paging.Pelem
}

func pager(ctx iface.Context, pagestr string, count, limit int) []paging.Pelem {
	if len(pagestr) == 0 {
		pagestr = "1"
	}
	if limit == 0 {
		return nil
	}
	p := ctx.NonPortable().Resource() + "?" + ctx.NonPortable().RawParams()
	page, err := strconv.Atoi(pagestr)
	if err != nil {
		return nil	// Not blowing up here.
	}
	if page == 0 {
		return nil
	}
	page_count := count/limit+1
	nav, _ := paging.P(page, page_count, 3, p)
	return nav
}

func elem(s interface{}, memb int64) interface{} {
	// Hotfix, we need to use reflection here.
	switch t := s.(type) {
	case []string:
		return t[memb]
	case []interface{}:
		return t[memb]
	}
	return "Error: unkown slice type."
}

func inSlice(s []interface{}, b interface{}) bool {
	for _, v := range s {
		if reflect.DeepEqual(v, b) {
			return true
		}
	}
	return false
}

type hook struct {
	ctx			iface.Context
	hook		iface.Hook
	hookName	string
}

func (h *hook) Has() bool {
	return h.hook.HasSubscribers()
}

func (h *hook) Fire(args ...interface{}) (string, error) {
	if !h.Has() {
		return "", nil
	}
	if h.hook.SubscriberCount() > 1 {
		return "", fmt.Errorf("More than one subscriber.")
	}
	h.hook.Fire(args...)
	subs := h.hook.Subscribers()
	return New(h.ctx).ToString([]string{subs[0].Name() + "/" + h.hookName})
}

func selectHook(ctx iface.Context, hookName string) *hook {
	return &hook{
		ctx,
		ctx.Conducting().Hooks().Select(hookName),
		hookName,
	}
}

// We must recreate this map each time because map write is not threadsafe.
// Write will happen when a hook modifies the map (hook call is not implemented yet).
func builtins(ctx iface.Context) map[string]interface{} {
	viewCtx := ctx.ViewContext().Get()
	user := ctx.User()
	ret := map[string]interface{}{
		"get": func(s ...string) interface{} {
			return get(viewCtx, s...)
		},
		"date": date,
		"is_stranger": func() bool {
			return user.Level() == 0
		},
		"logged_in": func() bool {
			return user.Level() > 0
		},
		"is_moderator": func() bool {
			return user.Level() >= 200
		},
		"is_admin": func() bool {
			return user.Level() >= 300
		},
		"is_map": isMap,
		"eq": eq,
		"html": html,
		"format_float": formatFloat,
		"fallback": fallback,
		"type_of":	typeOf,
		"same_kind": sameKind,
		"title": strings.Title,
		"url": func(action_name string, i ...interface{}) string {
			return _url(ctx, action_name, i...) 
		},
		"form": func(action_name string) *Form {
			return form(ctx, action_name)
		},
		"counter": newcounter,
		"get_sub": func(str string, params ...interface{}) []interface{} {
			return getSub(ctx, str, params...)
		},
		"get_list": func(str string, params ...interface{}) []interface{} {
			return getList(ctx, str, params...)
		},
		"elem": elem,
		"pager": func(pagesl []string, count, limited int) []paging.Pelem {
			var pagestr string
			if len(pagesl) == 0 {
				pagestr = "1"
			} else {
				pagestr = pagesl[0]
			}
			return pager(ctx, pagestr, count, limited)
		},
		"in_slice": inSlice,
		"hook": func(hookName string) *hook {
			return selectHook(ctx, hookName)
		},
	}
	ctx.Conducting().Hooks().Select("AddTemplateBuiltin").Fire(ret)
	return ret
}
