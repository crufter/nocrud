package terminal

import (
	"encoding/json"
	"fmt"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"reflect"
)

func latestOptdoc(db iface.Db) iface.Document {
	fil, err := db.NewFilter("options", map[string]interface{}{
		"sort": "-created",
	})
	if err != nil {
		panic(err)
	}
	o, err := fil.SelectOne()
	if err != nil {
		panic(err)
	}
	return o
}

func (c *C) install(resource, module string, ignore bool) error {
	modu := c.ctx.Conducting().Hooks().Module(module)
	if !modu.Exists() {
		return fmt.Errorf("Can't install nonexisting module \"%v\".", module)
	}
	odoc := latestOptdoc(c.ctx.Db())
	var err error
	if !ignore {
		upd := map[string]interface{}{
			"$addToSet": map[string]interface{}{
				fmt.Sprintf("nouns.%v.composedOf", resource): module,
			},
		}
		err = odoc.Update(upd)
		if err != nil {
			return err
		}
	}
	inst := modu.Instance()
	if inst.HasMethod("Install") {
		ret := func(r error) {
			err = r
		}
		inst.Method("Install").Call(ret, odoc, resource)
		if err != nil {
			return err
		}
	}
	if inst.HasMethod("InstallScript") {
		err = c.runScript(inst.Method("InstallScript"), resource)
	}
	return err
}

// Used when one wants to install a module without associating it to any resource.
func (c *C) installRaw(module string) error {
	return c.install("", module, true)
}

func (c *C) runScript(me iface.Method, resource string) error {
	var scriptStr string
	ret := func(s string) {
		scriptStr = s
	}
	me.Call(ret, resource)
	_, err := c.Execute(map[string]interface{}{
		"script": scriptStr,
	})
	return err
}

func (c *C) uninstall(resource, module string) error {
	modu := c.ctx.Conducting().Hooks().Module(module)
	if !modu.Exists() {
		return fmt.Errorf("Can't install nonexisting module \"%v\".", module)
	}
	odoc := latestOptdoc(c.ctx.Db())
	var err error
	if err != nil {
		return err
	}
	upd := map[string]interface{}{
		"$pull": map[string]interface{}{
			fmt.Sprintf("nouns.%v.composedOf", resource): module,
		},
	}
	err = odoc.Update(upd)
	if err != nil {
		return err
	}
	inst := modu.Instance()
	if inst.HasMethod("Uninstall") {
		ret := func(r error) {
			err = r
		}
		inst.Method("Uninstall").Call(ret, odoc, resource)
		if err != nil {
			return err
		}
	}
	if inst.HasMethod("UninstallScript") {
		err = c.runScript(inst.Method("UninstallScript"), resource)
	}
	return err
}

type Verb struct {
	Verb   string
	Module string
}

func (c *C) verbs(resource string) ([]Verb, error) {
	comp, ok := c.ctx.Options().Document().GetS(fmt.Sprintf("nouns.%v.composedOf", resource))
	if !ok {
		return nil, fmt.Errorf("Can't find resource %v.", resource)
	}
	verbs := map[string]Verb{}
	for _, v := range comp {
		modu := c.ctx.Conducting().Hooks().Module(v.(string))
		mtods := modu.Instance().MethodNames()
		for _, v1 := range mtods {
			if _, exists := verbs[v1]; exists {
				continue
			}
			verbs[v1] = Verb{
				v1,
				v.(string),
			}
		}
	}
	ret := []Verb{}
	for _, v := range verbs {
		ret = append(ret, v)
	}
	return ret, nil
}

func (c *C) composed(resource string) []string {
	return nil
}

func jsondec(s string) interface{} {
	return fmt.Errorf("Not implemented yet.")
}

func chost(s string) error {
	return fmt.Errorf("Not implemented yet.")
}

// Switches to a template.
// example usage:
// setTempl "private" "opesun"
func setTemplate(typ, name string) error {
	return fmt.Errorf("Not implemented yet.")
}

var argLabels = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

// Returns the help for a given command.
func help(funcMap map[string]interface{}, funcName string) string {
	fun, has := funcMap[funcName]
	if !has {
		return fmt.Sprintf("No function named %v.", funcName)
	}
	v := reflect.TypeOf(fun)
	ret := ""
	ret = ret + fmt.Sprintf("\nfunc %v(", funcName)
	for i := 0; i < v.NumIn(); i++ {
		ret = ret + argLabels[i] + " " + fmt.Sprint(v.In(i))
		if i < v.NumIn()-1 {
			ret = ret + ", "
		}
	}
	ret = ret + ") "
	if v.NumOut() > 1 {
		ret = ret + "("
	}
	for i := 0; i < v.NumOut(); i++ {
		ret = ret + fmt.Sprint(v.Out(i))
		if i < v.NumOut()-1 {
			ret = ret + ", "
		}
	}
	if v.NumOut() > 1 {
		ret = ret + ")"
	}
	return ret
}

// Returns true if the command available.
func avail(funcMap map[string]interface{}, funcName string) bool {
	_, has := funcMap[funcName]
	return has
}

// Returns the list of all commands.
func commands(b map[string]interface{}) []string {
	ret := []string{}
	for i := range b {
		ret = append(ret, i)
	}
	return ret
}

func tagOpt(s string) error {
	return fmt.Errorf("Not implemented yet.")
}

func saveOpt(s string) error {
	return fmt.Errorf("Not implemented yet.")
}

func revert(s string) error {
	return fmt.Errorf("Not implemented yet.")
}

func findCommand(s string) []string {
	return nil
}

func (c *C) setScheme(resource, jsn string) error {
	doc := latestOptdoc(c.ctx.Db())
	var v interface{}
	err := json.Unmarshal([]byte(jsn), &v)
	if err != nil {
		return fmt.Errorf("Json not valid:" + err.Error())
	}
	scheme := v.(map[string]interface{})
	upd := map[string]interface{}{
		"$set": map[string]interface{}{
			"nouns." + resource + ".verbs.Insert.input": scheme,
			"nouns." + resource + ".verbs.Update.input": scheme,
		},
	}
	return doc.Update(upd)
}

func unserVal(value interface{}) interface{} {
	if vStr, ok := value.(string); ok {
		var v interface{}
		err := json.Unmarshal([]byte(vStr), &v)
		if err == nil {
			return v
		}
	}
	return value
}

func (c *C) setNounAttr(noun, attr string, value interface{}) error {
	v := unserVal(value)
	o := latestOptdoc(c.ctx.Db())
	return o.Update(map[string]interface{}{
		"$set": map[string]interface{}{
			fmt.Sprintf("nouns.%v.%v", noun, attr): v,
		},
	})
}

func (c *C) setVerbAttr(noun, verb, attr string, value interface{}) error {
	v := unserVal(value)
	o := latestOptdoc(c.ctx.Db())
	return o.Update(map[string]interface{}{
		"$set": map[string]interface{}{
			fmt.Sprintf("nouns.%v.verb.%v.%v", noun, verb, attr): v,
		},
	})
}

func (c *C) setAttr(a string, value interface{}) error {
	v := unserVal(value)
	o := latestOptdoc(c.ctx.Db())
	return o.Update(map[string]interface{}{
		"$set": map[string]interface{}{
			a: v,
		},
	})
}

func (c *C) resetSite() error {
	o := latestOptdoc(c.ctx.Db())
	return o.Update(map[string]interface{}{})
}

func (c *C) builtins() map[string]interface{} {
	f := map[string]interface{}{
		"install": func(resource, module string) error {
			return c.install(resource, module, false)
		},
		"uninstall": func(resource, module string) error {
			return c.uninstall(resource, module)
		},
		"chost":    chost,
		"setTempl": setTemplate,
		"jsondec":  jsondec,
		"setScheme": func(resource, jsn string) error {
			return c.setScheme(resource, jsn)
		},
		"tagOpt":   tagOpt,
		"saveOpt":  saveOpt,
		"revert":   revert,
		"findComm": findCommand,
		"resetSite": func() error {
			return c.resetSite()
		},
		"verbs": func(s string) ([]Verb, error) {
			return c.verbs(s)
		},
		"composed": func(r string) []string {
			return c.composed(r)
		},
		"installRaw": func(module string) error {
			return c.installRaw(module)
		},
		"setNounAttr": func(a, b string, c1 interface{}) error {
			return c.setNounAttr(a, b, c1)
		},
		"setVerbAttr": func(a, b, c1 string, d interface{}) error {
			return c.setVerbAttr(a, b, c1, d)
		},
		"setAttr": func(a string, b interface{}) error {
			return c.setAttr(a, b)
		},
	}
	f["help"] = func(fname string) string {
		return help(f, fname)
	}
	f["avail"] = func(a string) bool {
		return avail(f, a)
	}
	f["commands"] = func() []string {
		return commands(f)
	}
	return f
}
