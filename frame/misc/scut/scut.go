// Package scut contains a somewhat ugly but useful collection of frequently appearing patterns to allow faster prototyping.
// Methods here are mainly related to view- or conroller-like parts.
package scut

import (
	"fmt"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"path/filepath"
	"strings"
)

// Gives you back the type of the currently used template (either "private" or public).
func TemplateType(opt map[string]interface{}) string {
	_, priv := opt["TplIsPrivate"]
	var ttype string
	if priv {
		ttype = "private"
	} else {
		ttype = "public"
	}
	return ttype
}

// Gives you back the name of the current template in use.
func TemplateName(opt map[string]interface{}) string {
	tpl, has_tpl := opt["Template"]
	if !has_tpl {
		tpl = "default"
	}
	return tpl.(string)
}

// Decides if a given relative filepath (filep) is a possible module filepath.
// This may be deprecated in the future since it seems so restrictive.
func PossibleModPath(filep string) bool {
	sl := strings.Split(filep, "/")
	return len(sl) >= 2
}

// TODO: Implement file caching here.
// Reads the filepath relative filepath from either the current template, or the fallback module tpl folder if filepath has at least one slash in it.
// file_reader is optional, falls back to simple ioutil.ReadFile if not given. file_reader will be a custom file_reader with caching soon.
func GetFile(fs iface.FileSys, filepath string) ([]byte, error) {
	templ, err := fs.SelectPlace("template")
	if err != nil {
		return nil, err
	}
	b, err := templ.File(filepath).Read()
	if err == nil {
		return b, nil
	}
	if !PossibleModPath(filepath) {
		return nil, fmt.Errorf("Not found.")
	}
	mod, err := fs.SelectPlace("modules")
	if err != nil {
		return nil, err
	}
	mtp := getModTPath(filepath)
	return mod.Directory(mtp.modname, "tpl").File(mtp.fpath).Read()
}

func Dirify(s string) string {
	return strings.Replace(s, ":", "-", -1)
}

// Observes opt and gives you back the path of your template eg
// "templates/public/template_name" or "templates/private/hostname/template_name"
func GetTPath(opt map[string]interface{}, host string) string {
	host = Dirify(host)
	templ := TemplateName(opt)
	ttype := TemplateType(opt)
	if ttype == "public" {
		return filepath.Join("templates", ttype, templ)
	}
	return filepath.Join("templates", ttype, host, templ)
}

type modTPath struct {
	modname		string
	fpath		string
}

// Inp:	"admin/this/that.txt"
// []string{ "modules/admin/tpl", "this/that.txt"}
func getModTPath(filename string) modTPath {
	if filename[0] == '/' {
		filename = filename[1:]
	}
	p := strings.Split(filename, "/")
	return modTPath{
		p[0],
		filepath.Join(p[1:]...),
	}
}

func NotAdmin(user iface.User) bool {
	return user.Level() < 300
}

func IsAdmin(user iface.User) bool {
	return user.Level() >= 300
}

func IsModerator(user iface.User) bool {
	return user.Level() >= 200
}

func IsRegistered(user iface.User) bool {
	return user.Level() >= 100
}

func IsStranger(user iface.User) bool {
	return user.Level() == 0
}

// Merges b into a (overwriting members in a.
func Merge(a map[string]interface{}, b map[string]interface{}) {
	for i, v := range b {
		a[i] = v
	}
}

// CanonicalHost(uni.Req.Host, uni.Opt)
// Gives you back the canonical address of the site so it can be made available from different domains.
func CanonicalHost(host string, opt map[string]interface{}) string {
	alias_whitelist, has_alias_whitelist := opt["host_alias_whitelist"]
	if has_alias_whitelist {
		awm := alias_whitelist.(map[string]interface{})
		if _, allowed := awm[host]; !allowed && len(awm) > 0 { // To prevent entirely locking yourself out of the site. Still can introduce problems if misused.
			panic(fmt.Sprintf("Unapproved host alias %v.", host))
		}
	}
	canon_host, has_canon := opt["canonical_host"]
	if !has_canon {
		return host
	}
	return canon_host.(string)
}

func OnlyAdmin(u iface.User) {
	if u.Level() < 300 {
		panic("Only an admin can do this operation.")
	}
}

// Gets the nouns from the options document.
// If the opt doc is empty, it returns a default nouns configoration.
func GetNouns(odoc iface.NestedData) map[string]interface{} {
	opt_def := map[string]interface{}{
		"composed_of": []interface{}{
			"jsonedit",
			"installer",
		},
	}
	nouns, ok := odoc.GetM("nouns")
	if !ok {
		nouns = map[string]interface{}{
			"options": opt_def,
		}
	}
	if _, ok := nouns["options"]; !ok {
		nouns["options"] = opt_def
	}
	return nouns
}