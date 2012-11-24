package terminal

import (
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/frame/misc/scut"
	"github.com/opesun/nocrud/modules/terminal/thack"
)

type C struct {
	ctx   iface.Context
	nouns map[string]interface{}
}

func (c *C) Init(ctx iface.Context) {
	c.ctx = ctx
	c.nouns = scut.GetNouns(ctx.Options().Document())
}

// Just to be able to load the Get template.
func (c *C) Get() {
}

func (c *C) Execute(data map[string]interface{}) (string, error) {
	script := data["script"].(string)
	return thack.New().Funcs(c.builtins()).Execute(script)
}

func (c *C) Install(o iface.Document, resource string) error {
	upd := map[string]interface{}{
		"$set": map[string]interface{}{
			"nouns." + resource + ".verbs.Execute.input": map[string]interface{}{
				"script": 1,
			},
		},
	}
	return o.Update(upd)
}

func (C *C) Uninstall(o iface.Document, resource string) error {
	upd := map[string]interface{}{
		"$set": map[string]interface{}{
			"nouns." + resource + ".verbs.Execute.input": 1,
		},
	}
	return o.Update(upd)
}
