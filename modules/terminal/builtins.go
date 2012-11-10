package terminal

import(
	iface "github.com/opesun/nocrud/frame/interfaces"
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

func (c *C) builtins() map[string]interface{} {
	return map[string]interface{}{
		"install": func(resource, module string) error {
			return c.install(resource, module)
		},
		"uninstall": func(resource, module string) error {
			return c.uninstall(resource, module)
		},
		"jsondec": jsondec,
	}
}