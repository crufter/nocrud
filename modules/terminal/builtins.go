package terminal

import (
	"encoding/json"
	"fmt"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"reflect"
)

func latestOptdoc(db iface.Db) (iface.Document, error) {
	fil, err := db.NewFilter("options", map[string]interface{}{
		"sort": "-created",
	})
	if err != nil {
		return nil, err
	}
	return fil.SelectOne()
}

func (c *C) install(resource, module string, ignore bool) error {
	modu := c.ctx.Conducting().Hooks().Module(module)
	if !modu.Exists() {
		return fmt.Errorf("Can't install nonexisting module \"%v\".", module)
	}
	odoc, err := latestOptdoc(c.ctx.Db())
	if err != nil {
		return err
	}
	if !ignore {
		upd := map[string]interface{}{
			"$addToSet": map[string]interface{}{
				fmt.Sprintf("nouns.%v.composed_of", resource): module,
			},
		}
		err = odoc.Update(upd)
		if err != nil {
			return err
		}
	}
	ret := func(r error) {
		err = r
	}
	inst := modu.Instance()
	if inst.HasMethod("Install") {
		inst.Method("Install").Call(ret, odoc, resource)
	}
	return err
}

// Used when one wants to install a module without associating it to any resource.
func (c *C) installRaw(module string) error {
	return c.install("", module, true)
}

func (c *C) uninstall(resource, module string) error {
	modu := c.ctx.Conducting().Hooks().Module(module)
	if !modu.Exists() {
		return fmt.Errorf("Can't install nonexisting module \"%v\".", module)
	}
	odoc, err := latestOptdoc(c.ctx.Db())
	if err != nil {
		return err
	}
	upd := map[string]interface{}{
		"$pull": map[string]interface{}{
			fmt.Sprintf("nouns.%v.composed_of", resource): module,
		},
	}
	err = odoc.Update(upd)
	if err != nil {
		return err
	}
	ret := func(r error) {
		err = r
	}
	inst := modu.Instance()
	if inst.HasMethod("Uninstall") {
		inst.Method("Uninstall").Call(ret, odoc, resource)
	}
	return err
}

type verb struct {
	Verb   string
	Module string
}

func (c *C) verbs(resource string) ([]verb, error) {
	comp, ok := c.ctx.Options().Document().GetS(fmt.Sprintf("nouns.%v.composed_of", resource))
	if !ok {
		return nil, fmt.Errorf("Can't find resource %v.", resource)
	}
	verbs := map[string]struct{}{}
	for _, v := range comp {
		modu := c.ctx.Conducting().Hooks().Module(v.(string))
		mtods := modu.Instance().MethodNames()
		fmt.Println(mtods)
		verbs[""] = struct{}{}
	}
	ret := []verb{}
	return ret, fmt.Errorf("Not implemented yet.")
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

// 
func (c *C) setScheme(resource, jsn string) error {
	doc, err := latestOptdoc(c.ctx.Db())
	if err != nil {
		return err
	}
	var v interface{}
	err = json.Unmarshal([]byte(jsn), &v)
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
		"composed": func(r string) []string {
			return c.composed(r)
		},
		"installRaw": func(module string) error {
			return c.installRaw(module)
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
