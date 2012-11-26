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
)

type m map[string]interface{}

type Top struct {
	ctx iface.Context
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

func (t *Top) route() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf(fmt.Sprint(r) + string(debug.Stack()))
		}
	}()
	ctx := t.ctx
	odoc := ctx.Options().Document()
	nouns := scut.GetNouns(odoc)
	np := ctx.NonPortable()
	hl, err := highlev.New(ctx.Conducting().Hooks(), np.Resource(), nouns, np.Params())
	if err != nil {
		display.New(ctx).Do(nil)
		return nil
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

func New(ctx iface.Context) *Top {
	return &Top{
		ctx,
	}
}
