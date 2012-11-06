package highlev

import(
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/frame/glue"
	"github.com/opesun/nocrud/frame/lang"
	"github.com/opesun/numcon"
	"github.com/opesun/sanitize"
	"github.com/opesun/jsonp"
	"fmt"
)

type HighLev struct {
	hooks		iface.Hooks
	resource	string
	nouns		map[string]interface{}
	params		map[string]interface{}
	desc		*glue.Descriptor
}

func New(hooks iface.Hooks, resource string, nouns, params map[string]interface{}) (*HighLev, error) {
	h := &HighLev{
		hooks:		hooks,
		resource: 	resource,
		nouns:		nouns,
		params:		params,
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

func (h *HighLev) Run(db iface.Db, usr iface.User, defminlev int) ([]interface{}, error) {
	desc := h.desc
	levi, ok := jsonp.Get(h.nouns, fmt.Sprintf("%v.verbs.%v.level", desc.Sentence.Noun, desc.Sentence.Verb))
	if !ok {
		levi = defminlev
	}
	lev, _ := numcon.Int(levi)
	if usr.Level() < lev {
		return nil, fmt.Errorf("Not allowed.")
	}
	filterCreator := func(c string, input map[string]interface{}) (iface.Filter, error) {
		return db.NewFilter(c, input)
	}
	inp, data, err := desc.CreateInputs(filterCreator)
	if err != nil {
		return nil, err
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
		params[form.KeyPrefix + i] = v
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