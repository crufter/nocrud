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

// Just to be able to load the Get template.
func (c *C) Get() {
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