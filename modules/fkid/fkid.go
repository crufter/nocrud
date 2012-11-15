package fkid

import(
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/sanitize"
	"strings"
)

type C struct {
	ctx		iface.Context
}

func (c *C) Init(ctx iface.Context) {
	c.ctx = ctx
}

func (c *C) SanitizerMangler(san *sanitize.Extractor) {
	san.AddFuncs(sanitize.FuncMap{
		"fkid": func(dat interface{}, s sanitize.Scheme) (interface{}, error) {
			cs, _ := s.Specific["comma_separated"].(bool)
			if cs {
				ret := []interface{}{}
				split := strings.Split(dat.(string), ",")
				for _, v := range split {
					idstr := strings.Trim(v, " ")
					id, err := c.ctx.Db().ToId(idstr)
					if err != nil {
						return nil, err
					}
					ret = append(ret, id)
				}
				return ret, nil
			}
			str := strings.Trim(dat.(string), " ")
			return c.ctx.Db().ToId(str)
		},
	})
}

func (c *C) FkidTypeHandler() {
}

func (c *C) Install(o iface.Document, resource string) error {
	upd := map[string]interface{}{
		"$addToSet": map[string]interface{}{
			"Hooks.fkidTypeHandler": "fkid",
		},
	}
	return o.Update(upd)
}

func (c *C) Uninstall(o iface.Document, resource string) error {
	upd := map[string]interface{}{
		"$pull": map[string]interface{}{
			"Hooks.fkidTypeHandler": "fkid",
		},
	}
	return o.Update(upd)
}