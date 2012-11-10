package terminal

import(
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/frame/misc/scut"
	"strings"
	"text/template"
	"bytes"
	"fmt"
)

type C struct {
	ctx 		iface.Context
	nouns 		map[string]interface{}
}

func (c *C) Init(ctx iface.Context) {
	c.ctx = ctx
	c.nouns = scut.GetNouns(ctx.Options().Document())
}

func latestOptdoc(db iface.Db) (iface.Document, error) {
	fil, err := db.NewFilter("options", map[string]interface{}{
		"sort": "-created",
	})
	if err != nil {
		return nil, err
	}
	return fil.SelectOne()
}

// Just to be able to load the Get template.
func (c *C) Get() {
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

func stripComments(lines []string) []string {
	ret := []string{}
	for _, v := range lines {
		if len(v) > 0 && v[0] != '#' {
			ret = append(ret, v)
		}
	}
	return ret
}

func strip(e error) error {
	return fmt.Errorf("line "+e.Error()[16:])
}

func (c *C) Execute(data map[string]interface{}) (string, error) {
	script := data["script"].(string)
	lines := strings.Split(script, "\n")
	for i, v := range lines {
		// This hack... :)
		if len(v) > 0 && v[0] != '#' {
			lines[i] = "{{" + v + "}}"
		}
	}
	whole_file := strings.Join(lines, "\n")
	funcMap := template.FuncMap(c.builtins())
	t, err := template.New("shell").Funcs(funcMap).Parse(string(whole_file))
	if err != nil {
		return "", strip(err)
	}
	context := map[string]interface{}{
	}
	var buffer bytes.Buffer
	err = t.Execute(&buffer, context)
	if err != nil {
		return "", strip(err)
	}
	output_lines := strings.Split(buffer.String(), "\n")
	output_lines = stripComments(output_lines)
	return strings.Join(output_lines, "\n"), nil
}