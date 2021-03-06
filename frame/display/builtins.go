package display

import (
	"fmt"
	"github.com/opesun/jsonp"
	"github.com/opesun/nocrud/frame/highlev"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/frame/lang"
	"github.com/opesun/nocrud/frame/misc/convert"
	"github.com/opesun/nocrud/frame/misc/scut"
	"github.com/opesun/numcon"
	"github.com/opesun/paging"
	"html/template"
	"reflect"
	"strconv"
	"strings"
	"time"
	"encoding/json"
)

func getMap(dat map[string]interface{}, s ...string) interface{} {
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
		ret = ret + `<input type="hidden" name="` + v[0] + `" value="` + v[1] + `" />`
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

func (c *counter) Inc() string { // Ugly hack, template engine needs a return value.
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
	return int(c)%i == 0
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
	deflev_i, _ := ctx.Options().Document().Get("defaultLevel")
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
	hl, err := highlev.New(ctx.Conducting().Hooks(), "/"+noun, nouns, inp)
	if err != nil {
		panic(err)
	}
	deflev_i, _ := ctx.Options().Document().Get("defaultLevel")
	deflev, _ := numcon.Int(deflev_i)
	ret, err := hl.Run(ctx.Db(), ctx.User(), deflev)
	if err != nil {
		panic(err)
	}
	return ret
}

type pagr struct {
	HasPrev bool
	Prev    int
	HasNext bool
	Next    int
	Elems   []paging.Pelem
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
		return nil // Not blowing up here.
	}
	if page == 0 {
		return nil
	}
	page_count := count/limit + 1
	nav, _ := paging.P(page, page_count, 3, p)
	return nav
}

func index(s interface{}, index int) interface{} {
	return reflect.ValueOf(s).Index(index).Interface()
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
	ctx      iface.Context
	hook     iface.Hook
	hookName string
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
	ma := convert.ListToMap(args...)
	inp := []interface{}{}
	for i, v := range ma {
		h.ctx.ViewContext().Publish(i, v)
		inp = append(inp, v)
	}
	subs := h.hook.Subscribers()
	modName := subs[0].Name()
	if h.ctx.Conducting().Hooks().Module(modName).Instance().HasMethod(h.hookName) {
		h.hook.Fire(inp...)
	}
	return New(h.ctx).ToString([]string{
		modName + "/" + strings.Title(h.hookName),
	})
}

func selectHook(ctx iface.Context, hookName string) *hook {
	return &hook{
		ctx,
		ctx.Conducting().Hooks().Select(hookName),
		hookName,
	}
}

func concat(s ...string) string {
	return strings.Join(s, "")
}

// Works with array, slice, map, chan, string.
func _len(a interface{}) int {
	return reflect.ValueOf(a).Len()
}

func setMap(a map[string]interface{}, b string, c interface{}) error {
	a[b] = c
	return nil	// We must have a return value according to html/template
}

func newMap() map[string]interface{} {
	return map[string]interface{}{}
}

func newSlice() []interface{} {
	return []interface{}{}
}

func indentedJSON(a interface{}) string {
	b, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(b)
}

// We must recreate this map each time because map write is not threadsafe.
// Write will happen when a hook modifies the map (hook call is not implemented yet).
func builtins(ctx iface.Context) map[string]interface{} {
	viewCtx := ctx.ViewContext().Get()
	user := ctx.User()
	ret := map[string]interface{}{
		"get": func(s ...string) interface{} {
			return getMap(viewCtx, s...)
		},
		"getMap": getMap,
		"date": date,
		"isStranger": func() bool {
			return user.Level() == 0
		},
		"loggedIn": func() bool {
			return user.Level() > 0
		},
		"isModerator": func() bool {
			return user.Level() >= 200
		},
		"isAdmin": func() bool {
			return user.Level() >= 300
		},
		"isMap":       	isMap,
		"eq":          	eq,
		"html":        	html,
		"formatFloat":	formatFloat,
		"newMap": 		newMap,
		"newSlice": 	newSlice,
		"fallback":    	fallback,
		"typeOf":      	typeOf,
		"sameKind":    	sameKind,
		"title":       	strings.Title,
		"url": func(action_name string, i ...interface{}) string {
			return _url(ctx, action_name, i...)
		},
		"form": func(action_name string) *Form {
			return form(ctx, action_name)
		},
		"counter": newcounter,
		"getSub": func(str string, params ...interface{}) []interface{} {
			return getSub(ctx, str, params...)
		},
		"getList": func(str string, params ...interface{}) []interface{} {
			return getList(ctx, str, params...)
		},
		"concat": concat,
		"index": index,
		"pager": func(page interface{}, count, limited int) []paging.Pelem {
			var pagestr string
			if page == nil {
				pagestr = "1"
			} else {
				pagestr = page.(string)
			}
			return pager(ctx, pagestr, count, limited)
		},
		"len": _len,
		"setMap": setMap,
		"inSlice": inSlice,
		"hook": func(hookName string) *hook {
			return selectHook(ctx, hookName)
		},
		"indentedJSON": indentedJSON,
	}
	ctx.Conducting().Hooks().Select("AddTemplateBuiltin").Fire(ret)
	return ret
}
