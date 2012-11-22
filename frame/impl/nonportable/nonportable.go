package nonportable

import (
	"net/http"
)

type NonPortable struct {
	path   string
	params map[string]interface{}
	req    *http.Request
	w      http.ResponseWriter
}

func New(path string, params map[string]interface{}, req *http.Request, w http.ResponseWriter) *NonPortable {
	return &NonPortable{
		path,
		params,
		req,
		w,
	}
}

func (n *NonPortable) Resource() string {
	return n.path
}

func (n *NonPortable) Params() map[string]interface{} {
	return n.params
}

func (n *NonPortable) ComingFrom() string {
	return n.req.Referer()
}

func (n *NonPortable) Redirect(redir string) {
	http.Redirect(n.w, n.req, redir, 303)
}

func (n *NonPortable) View() bool {
	return n.req.Method == "GET"
}

func (n *NonPortable) RawParams() string {
	return n.req.URL.RawQuery
}
