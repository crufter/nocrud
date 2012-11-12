package terminal

import(
	iface "github.com/opesun/nocrud/frame/interfaces"
	"reflect"
	"fmt"
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

func (c *C) install(resource, module string) error {
	modu := c.ctx.Conducting().Hooks().Module(module)
	if !modu.Exists() {
		return fmt.Errorf("Can't install nonexisting module \"%v\".", module)
	}
	odoc, err := latestOptdoc(c.ctx.Db())
	if err != nil {
		return err
	}
	upd := map[string]interface{}{
		"$addToSet": map[string]interface{}{
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
	if inst.HasMethod("Install") {
		inst.Method("Install").Call(ret, odoc, resource)
	}
	return err
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
	Verb	string
	Module	string
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

var arg_labels = []string{"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z"}

// Returns the help for a given command.
func help(funcmap map[string]interface{}, func_name string) string {
	fun, has := funcmap[func_name]
	if !has {
		return fmt.Sprintf("No function named %v.", func_name)
	}
	v := reflect.TypeOf(fun)
	ret := ""
	ret = ret + fmt.Sprintf("\nfunc %v(", func_name)
	for i:=0;i<v.NumIn();i++{
		ret = ret+arg_labels[i]+" "+fmt.Sprint(v.In(i))
		if i<v.NumIn()-1 {
			ret = ret + ", "
		}
	}
	ret = ret + ") "
	if v.NumOut() > 1 {
		ret = ret + "("
	}
	for i:=0;i<v.NumOut();i++{
		ret = ret + fmt.Sprint(v.Out(i))
		if i<v.NumOut()-1 {
			ret = ret + ", "
		}
	}
	if v.NumOut() > 1 {
		ret = ret + ")"
	}
	return ret
}

// Returns true if the command available.
func avail(funcmap map[string]interface{}, func_name string) bool {
	_, has := funcmap[func_name]
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

func (c *C) builtins() map[string]interface{} {
	f := map[string]interface{}{
		"install": func(resource, module string) error {
			return c.install(resource, module)
		},
		"uninstall": func(resource, module string) error {
			return c.uninstall(resource, module)
		},
		"chost": chost,
		"setTempl": setTemplate,
		"jsondec": jsondec,
		"tagOpt": tagOpt,
		"saveOpt": saveOpt,
		"revert": revert,
		"findComm": findCommand,
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