package bang

import (
	gt "github.com/opesun/gotrigga"
	"github.com/opesun/nocrud/frame/config"
	"github.com/opesun/nocrud/frame/impl/channels"
	"github.com/opesun/nocrud/frame/impl/client"
	"github.com/opesun/nocrud/frame/impl/conducting"
	"github.com/opesun/nocrud/frame/impl/context"
	"github.com/opesun/nocrud/frame/impl/db/mongodb"
	http_display "github.com/opesun/nocrud/frame/impl/display/http"
	ws_display "github.com/opesun/nocrud/frame/impl/display/ws"
	"github.com/opesun/nocrud/frame/impl/events"
	"github.com/opesun/nocrud/frame/impl/filesys"
	"github.com/opesun/nocrud/frame/impl/hooks"
	"github.com/opesun/nocrud/frame/impl/nonportable"
	"github.com/opesun/nocrud/frame/impl/options"
	"github.com/opesun/nocrud/frame/impl/temporaries"
	"github.com/opesun/nocrud/frame/impl/user"
	"github.com/opesun/nocrud/frame/impl/viewcontext"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/frame/misc/convert"
	"github.com/opesun/nocrud/frame/misc/scut"
	"github.com/opesun/nocrud/frame/mod"
	"io"
	"labix.org/v2/mgo"
	"mime/multipart"
	"net/http"
	"net/url"
)

func New(conn *gt.Connection, session *mgo.Session, dbConn *mgo.Database, w http.ResponseWriter, req *http.Request, config *config.Config) (iface.Context, error) {
	opts, _, err := queryOptions(dbConn, false)
	if err != nil {
		return nil, err
	}
	cli := client.New(w, req, req.Header, config.Secret)
	hookz, has := opts["Hooks"].(map[string]interface{})
	if !has {
		hookz = map[string]interface{}{}
	}
	hoo := hooks.New(hookz, mod.NewModule)
	datab := db.New(session, dbConn, opts, hoo)
	usr := user.New(datab, hoo, cli)
	err = req.ParseMultipartForm(1000000)
	if err != nil {
		return nil, err
	}
	var tempFiles map[string][]*multipart.FileHeader
	if req.MultipartForm != nil {
		tempFiles = req.MultipartForm.File
	}
	temp := temporaries.New(tempFiles)
	fs := filesys.New(config.AbsPath, scut.CanonicalHost(req.Host, opts), opts, temp)
	mods := modifiers(req.Form)
	np := nonportable.New(req.URL.Path, convert.Mapify(req.Form), req, w)
	vctx := viewcontext.New()
	dsp := http_display.New(w, config.AbsPath)
	ch := channels.New()
	o := options.New(opts, mods)
	ev := events.New(conn)
	cond := conducting.New(hoo, ev)
	ctx := context.New(cond, fs, usr, cli, datab, ch, vctx, np, dsp, o)
	initer := func(inst iface.Instance) error {
		if inst.HasMethod("Init") {
			inst.Method("Init").Call(nil, ctx)
		}
		return nil
	}
	hoo.Initer(initer)
	return ctx, err
}

func NewWS(conn *gt.Connection, session *mgo.Session, dbConn *mgo.Database, w http.ResponseWriter, req *http.Request, config *config.Config, ws io.Writer) (iface.Context, error) {
	opts, _, err := queryOptions(dbConn, false)
	if err != nil {
		return nil, err
	}
	cli := client.New(w, req, req.Header, config.Secret)
	hookz, has := opts["Hooks"].(map[string]interface{})
	if !has {
		hookz = map[string]interface{}{}
	}
	hoo := hooks.New(hookz, mod.NewModule)
	datab := db.New(session, dbConn, opts, hoo)
	usr := user.New(datab, hoo, cli)
	err = req.ParseMultipartForm(1000000)
	if err != nil {
		return nil, err
	}
	var tempFiles map[string][]*multipart.FileHeader
	if req.MultipartForm != nil {
		tempFiles = req.MultipartForm.File
	}
	temp := temporaries.New(tempFiles)
	fs := filesys.New(config.AbsPath, scut.CanonicalHost(req.Host, opts), opts, temp)
	mods := modifiers(req.Form)
	np := nonportable.New(req.URL.Path, convert.Mapify(req.Form), req, w)
	vctx := viewcontext.New()
	dsp := ws_display.New(ws, config.AbsPath)
	ch := channels.New()
	o := options.New(opts, mods)
	ev := events.New(conn)
	cond := conducting.New(hoo, ev)
	ctx := context.New(cond, fs, usr, cli, datab, ch, vctx, np, dsp, o)
	initer := func(inst iface.Instance) error {
		if inst.HasMethod("Init") {
			inst.Method("Init").Call(nil, ctx)
		}
		return nil
	}
	hoo.Initer(initer)
	return ctx, err
}

// Strips information unrelated to verb input from the Form.
func modifiers(a url.Values) map[string]interface{} {
	flags := []string{"json", "src", "nofmt", "ok", "action", "redirect"}
	mods := map[string]interface{}{}
	for _, v := range flags {
		if val, has := a[v]; has {
			mods[v] = val
			delete(a, v)
		}
	}
	for i, v := range a {
		if i[0] == '-' {
			mods[i[1:]] = v
			delete(a, i)
		}
	}
	return mods
}
