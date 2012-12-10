package highlev

import (
	"fmt"
	"github.com/opesun/jsonp"
	"github.com/opesun/nocrud/frame/glue"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/frame/lang"
	"github.com/opesun/numcon"
	"github.com/opesun/sanitize"
)

type HighLev struct {
	hooks    iface.Hooks
	resource string
	nouns    map[string]interface{}
	params   map[string]interface{}
	desc     *glue.Descriptor
}

func New(hooks iface.Hooks, resource string, nouns, params map[string]interface{}) (*HighLev, error) {
	h := &HighLev{
		hooks:    hooks,
		resource: resource,
		nouns:    nouns,
		params:   params,
	}
	desc, err := h.createDesc()
	if err != nil {
		return nil, err
	}
	h.desc = desc
	return h, nil
}

func (h *HighLev) createDesc() (*glue.Descriptor, error) {
	desc, err := glue.Identify(h.resource, h.nouns, h.params)
	if err != nil {
		return nil, err
	}
	return desc, nil
}

func (h *HighLev) userChecks(user map[string]interface{}) error {
	verbSpec, ok := jsonp.GetM(h.nouns, fmt.Sprint("%v.verbs.%v.userCrit", h.desc.Sentence.Noun, h.desc.Sentence.Verb))
	if ok {
		_, err := sanitize.Fast(verbSpec, user)
		return err
	}
	nounSpec, ok := jsonp.GetM(h.nouns, fmt.Sprintf("%v.userCrit", h.desc.Sentence.Noun))
	if !ok {
		return nil
	}
	_, err := sanitize.Fast(nounSpec, user)
	return err
}

func (h *HighLev) ownLev() (int, bool) {
	verbL, ok := jsonp.Get(h.nouns, fmt.Sprintf("%v.verbs.%v.ownLevel", h.desc.Sentence.Noun, h.desc.Sentence.Verb))
	if ok {
		return numcon.IntP(verbL), true
	}
	nounL, ok := jsonp.GetM(h.nouns, fmt.Sprintf("%v.ownLevel", h.desc.Sentence.Noun))
	if ok {
		return numcon.IntP(nounL), true
	}
	return 0, false
}

func (h *HighLev) Run(db iface.Db, usr iface.User, defminlev int) ([]interface{}, error) {
	desc := h.desc
	// Authentication.
	levi, ok := jsonp.Get(h.nouns, fmt.Sprintf("%v.verbs.%v.level", desc.Sentence.Noun, desc.Sentence.Verb))
	if !ok {
		levi = defminlev
	}
	lev, _ := numcon.Int(levi)
	if usr.Level() < lev {
		return nil, fmt.Errorf("Not allowed.")
	}
	err := h.userChecks(usr.Data())
	if err != nil {
		return nil, err
	}
	filterCreator := func(c string, input map[string]interface{}) (iface.Filter, error) {
		return db.NewFilter(c, input)
	}
	inp, data, err := desc.CreateInputs(filterCreator)
	if err != nil {
		return nil, err
	}
	ownLev, own := h.ownLev()
	if len(inp) > 0 {
		if f, ok := inp[0].(iface.Filter); ok {
			// This hook allows you to modify a filter before a verb accesses it.
			h.hooks.Select("TopModFilter").Fire(f)
			h.hooks.Select(f.Subject() + "TopModFilter").Fire(f)
		}
		if own && lev <= ownLev {
			if f, ok := inp[0].(iface.Filter); ok {
				f.AddQuery(map[string]interface{}{
					"createdBy": usr.Id(),
				})
			}
		}
	}
	if data != nil {
		if desc.Sentence.Noun != "options" {
			data, err = h.validate(desc.Sentence.Noun, desc.Sentence.Verb, data)
			if err != nil {
				return nil, err
			}
		}
		inp = append(inp, data)
	}
	module := h.hooks.Module(desc.VerbLocation)
	if !module.Exists() {
		return nil, fmt.Errorf("Unkown module.")
	}
	ins := module.Instance()
	var ret []interface{}
	ret_rec := func(i ...interface{}) {
		ret = i
	}
	err = ins.Method(desc.Sentence.Verb).Call(ret_rec, inp...)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (h *HighLev) Noun() string {
	return h.desc.Sentence.Noun
}

func (h *HighLev) Verb() string {
	return h.desc.Sentence.Verb
}

func (h *HighLev) VerbLocation() string {
	return h.desc.VerbLocation
}

func (h *HighLev) URLE() *lang.URLEncoder {
	return lang.NewURLEncoder(h.desc.Route, h.desc.Sentence)
}

func (h *HighLev) Sub(actionOrNoun string, p map[string]interface{}) (*HighLev, error) {
	form := lang.NewURLEncoder(h.desc.Route, h.desc.Sentence).Form(actionOrNoun)
	params := form.FilterFields
	for i, v := range p {
		params[form.KeyPrefix+i] = v
	}
	return New(h.hooks, form.ActionPath, h.nouns, params)
}

func (h *HighLev) validate(noun, verb string, data map[string]interface{}) (map[string]interface{}, error) {
	scheme_map, ok := jsonp.GetM(h.nouns, fmt.Sprintf("%v.verbs.%v.input", noun, verb))
	if !ok {
		return nil, fmt.Errorf("Can't find scheme for %v %v.", noun, verb)
	}
	ex, err := sanitize.New(scheme_map)
	if err != nil {
		return nil, err
	}
	h.hooks.Select("SanitizerMangler").Fire(ex)
	data, err = ex.Extract(data)
	if err != nil {
		return nil, err
	}
	h.hooks.Select("SanitizedDataMangler").Fire(data)
	return data, nil
}
