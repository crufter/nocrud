// Thack is a script language-like thingie hacked out of the golang template package.
package thack

import(
	"text/template"
	"strings"
	"bytes"
	"fmt"
)

type Thack struct {
	funcMap map[string]interface{}
}

func New() *Thack {
	return &Thack{
		map[string]interface{}{},
	}
}

// Not threadsafe.
func (t *Thack) Funcs(f map[string]interface{}) *Thack {
	t.funcMap = f
	return t
}

func (t *Thack) Execute(s string) (string, error) {
	lines := strings.Split(s, "\n")
	for i, v := range lines {
		// This hack... :)
		if len(v) > 0 && v[0] != '#' {
			lines[i] = "{{" + v + "}}"
		}
	}
	whole_file := strings.Join(lines, "\n")
	templ := template.New("shell")
	if t.funcMap != nil {
		templ.Funcs(t.funcMap)
	}
	templ, err := templ.Parse(string(whole_file))
	if err != nil {
		return "", strip(err)
	}
	context := map[string]interface{}{}
	var buffer bytes.Buffer
	err = templ.Execute(&buffer, context)
	if err != nil {
		return "", strip(err)
	}
	output_lines := strings.Split(buffer.String(), "\n")
	output_lines = stripComments(output_lines)
	return strings.Join(output_lines, "\n"), nil
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
	return fmt.Errorf("line " + e.Error()[16:])
}