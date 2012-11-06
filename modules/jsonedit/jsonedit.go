package jsonedit

import(
	"encoding/json"
	"github.com/opesun/nocrud/frame/composables/basics"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"time"
	"fmt"
)

type C struct {
			basics.Basics
	opt		map[string]interface{}
}

func (c *C) Init(ctx iface.Context) {
	c.Basics.Hooks = ctx.Conducting().Hooks()
	c.Basics.Db = ctx.Db()
}

func (c *C) decrypt(data map[string]interface{}) (map[string]interface{}, error) {
	jsond, ok := data["json"].(string)
	if !ok {
		return nil, fmt.Errorf("Member json is nonexistent or not a string.")
	}
	var v interface{}
	err := json.Unmarshal([]byte(jsond), &v)
	if err != nil {
		return nil, err
	}
	return v.(map[string]interface{}), nil
}

func (c *C) Insert(a iface.Filter, data map[string]interface{}) error {
	m, err := c.decrypt(data)
	if err != nil {
		return err
	}
	m["created"] = time.Now().UnixNano()	// Should include user too maybe.
	return a.Insert(m)
}

func (c *C) Update(a iface.Filter, data map[string]interface{}) error {
	m, err := c.decrypt(data)
	if err != nil {
		return err
	}
	m["modified"] = time.Now().UnixNano()
	return a.Update(m)
}

func (c *C) New() error {
	return nil
}

var ignore = []string{"_id", "created", "modified"}

func (c *C) Edit(a iface.Filter) (string, error) {
	doc, err := a.FindOne()
	if err != nil {
		return "", err
	}
	for _, v := range ignore {
		delete(doc, v)
	}
	marsh, err := json.MarshalIndent(doc, "", "\t")
	if err != nil {
		return "", err
	}
	return string(marsh), nil
}