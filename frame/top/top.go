package top

import (
	"fmt"
	"github.com/opesun/nocrud/frame/display"
	"github.com/opesun/nocrud/frame/highlev"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/frame/misc/scut"
	"github.com/opesun/nocrud/frame/verbinfo"
	"github.com/opesun/numcon"
	"runtime/debug"
	"strconv"
	"strings"
)

type m map[string]interface{}

type Top struct {
	ctx iface.Context
}

func New(ctx iface.Context) *Top {
	return &Top{
		ctx,
	}
}

func burnResults(vctx iface.ViewContext, key string, b []interface{}) {
	for i, v := range b {
		if i == 0 {
			vctx.Publish(key, v)
		} else {
			vctx.Publish(key+strconv.Itoa(i), v)
		}
	}
}

func (t *Top) Get(ret []interface{}, files []string) error {
	ran := verbinfo.NewRanalyzer(ret)
	if ran.HadError() {
		return display.New(t.ctx).Error(ran.Error())
	}
	if ran.ReturnedFile() {
		return t.ctx.NonPortable().ServeFile(ran.File())
	}
	burnResults(t.ctx.ViewContext(), "main", ran.NonErrors())
	return display.New(t.ctx).Do(files)
}

func (t *Top) Post(ret []interface{}, verb string) error {
	ran := verbinfo.NewRanalyzer(ret)
	var err error
	if ran.HadError() {
		err = ran.Error()
	}
	t.actionResponse(err, verb, ran.NonErrors())
	return nil
}

func (t *Top) Route() error {
	var err error
	err = t.route()
	if err != nil {
		t.ctx.Display().Write([]byte(err.Error()))
	}
	return nil
}

func (t *Top) RouteWS() error {
	err := t.routeWS()
	if err != nil {
		t.ctx.Display().Write([]byte(err.Error()))
	}
	return nil
}

// Files can not be accessed here.
func (t *Top) routeWS() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf(fmt.Sprint(r) + string(debug.Stack()))
		}
	}()
	ctx := t.ctx
	odoc := ctx.Options().Document()
	nouns, ok := odoc.GetM("nouns")
	if !ok {
		return fmt.Errorf("No nouns, cant route websocket.")
	}
	np := ctx.NonPortable()
	hl, err := highlev.New(ctx.Conducting().Hooks(), np.Resource(), nouns, np.Params())
	if err != nil {
		return err
	}
	deflev_i, _ := odoc.Get("defaultLevel")
	deflev, _ := numcon.Int(deflev_i)
	_, err = hl.Run(ctx.Db(), ctx.User(), deflev)
	return err
}

// Returns true if the path s identifies a file.
func isFile(s string) bool {
	sl := strings.Split(s, "/")
	return strings.Index(sl[len(sl)-1], ".") != -1
}

func (t *Top) route() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf(fmt.Sprint(r) + string(debug.Stack()))
		}
	}()
	ctx := t.ctx
	if isFile(ctx.NonPortable().Resource()) {
		return t.serveFile()
	}
	odoc := ctx.Options().Document()
	nouns := scut.GetNouns(odoc)
	np := ctx.NonPortable()
	hl, err := highlev.New(ctx.Conducting().Hooks(), np.Resource(), nouns, np.Params())
	if err != nil {
		return display.New(ctx).Do(nil)
	}
	deflev_i, _ := odoc.Get("defaultLevel")
	deflev, _ := numcon.Int(deflev_i)
	ret, err := hl.Run(ctx.Db(), ctx.User(), deflev)
	if err != nil {
		return err
	}
	if ctx.NonPortable().View() {
		ctx.ViewContext().Publish("main_noun", hl.Noun())
		files := []string{
			hl.Noun() + "/" + hl.Verb(),
			hl.VerbLocation() + "/" + hl.Verb(),
		}
		return t.Get(ret, files)
	} else {
		return t.Post(ret, hl.Verb())
	}
	return nil
}

// Since we don't include the template name into the url, only "template", we have to extract the template name from the opt here.
// Example: xyz.com/template/style.css
//			xyz.com/tpl/admin/style.css
func (t *Top) serveFile() error {
	path := t.ctx.NonPortable().Resource()
	parts := strings.Split(path, "/")
	firstPart := parts[1]	// Intentionally ignoring the 0th one here.
	lastPart := parts[len(parts)-1]
	isGoFile := strings.HasSuffix(lastPart, ".go")
	if isGoFile {
		return fmt.Errorf("Can't serve source.")
	}
	if firstPart == "template" || firstPart == "tpl" {
		return t.serveTemplateFile()
	} else if firstPart == "uploads" {
		return t.serveUploadedFile()
	}
	rootPlace, err := t.ctx.FileSys().SelectPlace("root")
	if err != nil {
		return err
	}
	fileToServe := rootPlace.File(path)
	return t.ctx.NonPortable().ServeFile(fileToServe)
}

func (t *Top) serveTemplateFile() error {
	path := t.ctx.NonPortable().Resource()
	parts := strings.Split(path, "/")
	if parts[1] == "template" {
		templatePlace, err := t.ctx.FileSys().SelectPlace("template")
		if err != nil {
			return err
		}
		fileToServe := templatePlace.File(parts[2:]...)
		return t.ctx.NonPortable().ServeFile(fileToServe)
	}
	// "tpl"
	modulesPlace, err := t.ctx.FileSys().SelectPlace("modules")
	if err != nil {
		return err
	}
	a := []string{parts[2], "tpl"}
	a = append(a, parts[3:]...)
	fileToServe := modulesPlace.File(a...)
	return t.ctx.NonPortable().ServeFile(fileToServe)
}

func (t *Top) serveUploadedFile() error {
	path := t.ctx.NonPortable().Resource()
	parts := strings.Split(path, "/")
	uploadsPlace, err := t.ctx.FileSys().SelectPlace("uploads")
	if err != nil {
		return err
	}
	fileToServe := uploadsPlace.File(parts[2:]...)
	return t.ctx.NonPortable().ServeFile(fileToServe)
}

