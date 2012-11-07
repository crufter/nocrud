package bang

import(
	"labix.org/v2/mgo"
	"net/http"
	"net/url"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/frame/config"
	"github.com/opesun/nocrud/frame/mod"
	"github.com/opesun/nocrud/frame/misc/convert"
	"github.com/opesun/nocrud/frame/misc/scut"
	"github.com/opesun/nocrud/frame/impl/filter"
	"github.com/opesun/nocrud/frame/impl/hooks"
	"github.com/opesun/nocrud/frame/impl/user"
	"github.com/opesun/nocrud/frame/impl/db/mongodb"
	"github.com/opesun/nocrud/frame/impl/client"
	"github.com/opesun/nocrud/frame/impl/events"
	"github.com/opesun/nocrud/frame/impl/set/mongodb"
	"github.com/opesun/nocrud/frame/impl/temporaries"
	"github.com/opesun/nocrud/frame/impl/filesys"
	"github.com/opesun/nocrud/frame/impl/nonportable"
	"github.com/opesun/nocrud/frame/impl/viewcontext"
	http_display "github.com/opesun/nocrud/frame/impl/display/http"
	ws_display "github.com/opesun/nocrud/frame/impl/display/ws"
	gt "github.com/opesun/gotrigga"
	"github.com/opesun/nocrud/frame/impl/channels"
	"github.com/opesun/nocrud/frame/impl/conducting"
	"github.com/opesun/nocrud/frame/impl/options"
	"github.com/opesun/nocrud/frame/impl/context"
	"path/filepath"
	"strings"
	"io"
	"mime/multipart"
)

func New(conn *gt.Connection, session *mgo.Session, dbConn *mgo.Database, w http.ResponseWriter, req *http.Request, config *config.Config) (iface.Context, error) {
	opts, _, err := queryOptions(dbConn, false)
	if err != nil {
		return nil, err
	}
	paths := strings.Split(req.URL.Path, "/")
	if strings.Index(paths[len(paths)-1], ".") != -1 {
		serveFile(opts, w, req, config.AbsPath, scut.CanonicalHost(req.Host, opts), req.URL.Path)
		return nil, nil
	}
	cli := client.New(w, req, req.Header, config.Secret)
	hookz, has := opts["Hooks"].(map[string]interface{})
	if !has {
		hookz = map[string]interface{}{}
	}
	hoo := hooks.New(hookz, mod.NewModule)
	datab := db.New(session, dbConn, opts, hoo)
	usrFilter, err := filter.NewSimple(set.New(dbConn, "users"), nil)
	if err != nil {
		return nil, err
	}
	usr := user.New(datab, usrFilter, hoo, cli)
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
	usrFilter, err := filter.NewSimple(set.New(dbConn, "users"), nil)
	if err != nil {
		return nil, err
	}
	usr := user.New(datab, usrFilter, hoo, cli)
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

// File serving is migrated here temporarily.

// Since we don't include the template name into the url, only "template", we have to extract the template name from the opt here.
// Example: xyz.com/template/style.css
//			xyz.com/tpl/admin/style.css
func serveFile(opt map[string]interface{}, w http.ResponseWriter, req *http.Request, absPath, host string, path string) {
	paths := strings.Split(path, "/")
	first_p := paths[1]
	last_p := paths[len(paths)-1]
	has_sfx := strings.HasSuffix(last_p, ".go")
	if first_p == "template" || first_p == "tpl" && !has_sfx {
		serveTemplateFile(opt, w, req, absPath, host, path)
	} else if first_p == "uploads" {
		serveUploadedFile(opt, w, req, absPath, host, path)
	} else if !has_sfx {
		if paths[1] == "shared" {
			http.ServeFile(w, req, filepath.Join(absPath, path))
		} else {
			http.ServeFile(w, req, filepath.Join(absPath, "uploads", host, path))
		}
	}
}

func serveTemplateFile(opt map[string]interface{}, w http.ResponseWriter, req *http.Request, absPath, host string, path string) {
	paths := strings.Split(path, "/")
	if paths[1] == "template" {
		p := scut.GetTPath(opt, host)
		http.ServeFile(w, req, filepath.Join(absPath, p, strings.Join(paths[2:], "/")))
	} else { // "tpl"
		http.ServeFile(w, req, filepath.Join(absPath, "modules", paths[2], "tpl", strings.Join(paths[3:], "/")))
	}
}

func serveUploadedFile(opt map[string]interface{}, w http.ResponseWriter, req *http.Request, absPath, host string, path string) {
	paths := strings.Split(path, "/")
	http.ServeFile(w, req, filepath.Join(absPath, "uploads", scut.Dirify(host), strings.Join(paths[2:], "/")))
}