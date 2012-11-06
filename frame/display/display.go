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
	d.ctx.ViewContext().Publish("form", d.ctx.NonPortable().Params())
	beforeDisplay(d.ctx.Conducting().Hooks())
	loc, err := display_model.LoadLocStrings(d.ctx.FileSys(), d.ctx.ViewContext(), d.ctx.User().Languages()) // TODO: think about errors here.
	if err != nil {
		return err
	}
	d.ctx.ViewContext().Publish("loc", loc)
	if d.ctx.Options().Modifiers().Exists("json") {
		d.putJSON()
		return nil
	}
	return d.file(point)
}

func beforeDisplay(hooks iface.Hooks) {
	defer func(){
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()
	hooks.Select("BeforeDisplay").Fire()
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
	//file = append([]byte(fmt.Sprintf("<!-- %v/%v. -->", root, fi)), file...)
	//file = append(file, []byte(fmt.Sprintf("<!-- /%v/%v -->", root, fi))...)
	return file, nil
}

// Tries to dislay a template file.
func (d *Display) file(filepath string) error {
	src := false
	fileContents, err := require.R("", filepath + ".tpl",
	func(root, fi string) ([]byte, error) {
		return getFileAndConvert(d.ctx.FileSys(), fi)
	})
	if err != nil {
		return fmt.Errorf("Cant find template file %v.", filepath)
	}
	if src {
		displ:= d.ctx.Display()
		displ.Type("html").Write([]byte(fileContents))
	}
	err = d.prepareAndExec(fileContents)
	if err != nil {
		panic(err)
	}
	return nil
}

// Loads localization, template functions and executes the template.
func (d *Display) prepareAndExec(fileContents string) error {
	loc, err := display_model.LoadLocTempl(d.ctx.FileSys(), fileContents, d.ctx.User().Languages()) // TODO: think about errors here.
	if err != nil {
		return err
	}
	d.ctx.ViewContext().Publish("loc", loc)
	funcMap := template.FuncMap(builtins(d.ctx))
	t, err := template.New("tpl").Funcs(funcMap).Parse(fileContents)
	if err != nil {
		return err
	}
	displ := d.ctx.Display()
	displ.Type("html")
	return t.Execute(displ.Writer(), d.ctx.ViewContext().Get()) // TODO: watch for errors in execution.
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
	d.ctx.ViewContext().Publish("error", err)
	if d.ctx.Options().Modifiers().Exists("json") {
		return d.putJSON()
	}
	return d.file("error")
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
