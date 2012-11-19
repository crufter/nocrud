package display

import (
	"encoding/json"
	"fmt"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/frame/misc/scut"
	"github.com/opesun/nocrud/frame/display/model"
	"github.com/opesun/jsonp"
	"github.com/opesun/require"
	"github.com/russross/blackfriday"
	"html/template"
	"strings"
	"io"
	"bytes"
)

type Display struct {
	ctx		iface.Context
}

func New(ctx iface.Context) *Display {
	return &Display {
		ctx: ctx,
	}
}

func (d *Display) Do(files []string) error {
	point := d.decidePoint(files)
	d.publishForm()
	runDisplayHook(d.ctx.Conducting().Hooks())
	err := d.localizeContext()
	if err != nil {
		return err
	}
	displ := d.ctx.Display()
	modif := d.ctx.Options().Modifiers()
	if modif.Exists("json") {
		return d.putJSON()
	}
	file, err := d.getFile(point)
	if err != nil {
		return err
	}
	if modif.Exists("src") {
		return displ.Type("html").Write([]byte(file))
	}
	if d.ctx.Options().Modifiers().Exists("json") {
		return d.putJSON()
	}
	return d.execute(displ.Writer(), file)
}

func (d *Display) publishForm() {
	d.ctx.ViewContext().Publish("form", d.ctx.NonPortable().Params())
}

func (d *Display) localizeContext() error {
	loc, err := display_model.LoadLocStrings(d.ctx.FileSys(), d.ctx.ViewContext(), d.ctx.User().Languages()) // TODO: think about errors here.
	if err != nil {
		return err
	}
	d.ctx.ViewContext().Publish("loc", loc)
	return nil
}

func runDisplayHook(hooks iface.Hooks) {
	defer func(){
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()
	hooks.Select("BeforeDisplay").Fire()
}

func (d *Display) ToString(files []string) (string, error) {
	point := existing(d.ctx.FileSys(), files)
	file, err := d.getFile(point)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	err = d.execute(buf, file)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (d *Display) decidePoint(files []string) string {
	var point string
	if len(files) > 0 {
		point = existing(d.ctx.FileSys(), files)
	} else {
		p := d.ctx.NonPortable().Resource()
		if p == "/" {
			point = "index"
		} else {
			point = p
		}
	}
	return point
}

func merge(a interface{}, b map[string]interface{}) map[string]interface{} {
	if a == nil {
		return b
	}
	if b == nil {
		return a.(map[string]interface{})
	}
	a_m := a.(map[string]interface{})
	for i, v := range b {
		a_m[i] = v
	}
	return a_m
}

func toStringSlice(a interface{}) []string {
	if a == nil {
		return nil
	}
	switch val := a.(type) {
	case []interface{}:
		return jsonp.ToStringSlice(val)
	case []string:
		return val
	}
	return nil
}

func validFormat(format string) bool {
	switch format {
	case "md":
		return true
	}
	return false
}

func (d *Display) publishLangs() {
	langs := []interface{}{}
	for _, v := range d.ctx.User().Languages() {
		langs = append(langs, v)
	}
	d.ctx.ViewContext().Publish("langs", langs)
}

// Tries to dislay a template file.
func (d *Display) getFile(filepath string) (string, error) {
	return require.R("", filepath + ".tpl",
	func(root, fi string) ([]byte, error) {
		return getFileAndConvert(d.ctx.FileSys(), fi)
	})
}

// Loads localization, template functions and executes the template.
func (d *Display) execute(wr io.Writer, fileContents string) error {
	loc, err := display_model.LoadLocTempl(d.ctx.FileSys(), fileContents, d.ctx.User().Languages()) // TODO: think about errors here.
	if err != nil {
		return err
	}
	vctx := d.ctx.ViewContext()
	vctx.Publish("loc", loc)
	funcMap := template.FuncMap(builtins(d.ctx))
	t, err := template.New("tpl").Funcs(funcMap).Parse(fileContents)
	if err != nil {
		return err
	}
	return t.Execute(wr, vctx.Get()) // TODO: watch for errors in execution.
}

// Prints all available data to http response as a JSON.
func (d *Display) putJSON() error {
	var v []byte
	var err error
	nofmt := false
	if nofmt {
		v, err = json.Marshal(d.ctx.ViewContext().Get())
	} else {
		v, err = json.MarshalIndent(d.ctx.ViewContext().Get(), "", "    ")
	}
	if err != nil {
		return err
	}
	displ := d.ctx.Display()
	return displ.Type("json").Write(v)
}

// This is called if an error occured in a front hook.
func (d *Display) Error(err error) error {
	d.ctx.ViewContext().Publish("error", err.Error())
	d.Do([]string{"error"})
	return err	// Intentionally not returning the error coming from the above line.
}

// Does format conversions.
// Currently only: markdown -> html
func getFileAndConvert(fs iface.FileSys, filepath string) ([]byte, error) {
	file, err := scut.GetFile(fs, filepath)
	if err != nil {
		return nil, err
	}
	spl := strings.Split(filepath, ".")
	extension := spl[len(spl)-1]
	// In tpl files the first line contains the extension information, like "--md". (An entry point can't change it's extension.)
	if extension == "tpl" {
		strfile := string(file)
		newline_pos := strings.Index(strfile, "\n")
		if newline_pos > 3 && validFormat(strfile[2:newline_pos-1]) { // "--" plus at least 1 characer.
			extension = strfile[2 : newline_pos-1]
			file = file[newline_pos:]
		}
	}
	switch extension {
	case "md":
		file = blackfriday.MarkdownCommon(file)
	}
	return file, nil
}

func existing(f iface.FileSys, s []string) string {
	if len(s) > 10 {
		panic("Ouch.")
	}
	templ, err := f.SelectPlace("template")
	if err != nil {
		panic(err)
	}
	for i, v := range s {
		if i == len(s)-1 {
			return v
		}
		ex, err := templ.File(v + ".tpl").Exists()
		if err != nil {
			panic(err)
		}
		if ex {
			return v
		}
	}
	return ""
}
