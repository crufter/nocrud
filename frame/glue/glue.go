package glue

import(
	"github.com/opesun/nocrud/frame/mod"
	"github.com/opesun/nocrud/frame/lang"
	"github.com/opesun/nocrud/frame/lang/speaker"
	"fmt"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/frame/verbinfo"
)

type Descriptor struct {
	Route 			*lang.Route
	Sentence 		*lang.Sentence
	VerbLocation	string						// Name of the module with the verb.
	nounOpt			map[string]interface{}
}

func moduleHasVerb(modname string, verbname string) bool {
	mo := mod.NewModule(modname)
	if !mo.Exists() {
		return false
	}
	return mo.Instance().HasMethod(verbname)
}

func Identify(path string, nouns map[string]interface{}, inp map[string]interface{}) (*Descriptor, error) {
	desc := &Descriptor{}
	r, err := lang.NewRoute(path, inp)
	if err != nil {
		return nil, err
	}
	desc.Route = r
	speaker := speaker.New(moduleHasVerb, nouns)
	sentence, err := lang.NewSentence(r, speaker)
	if err != nil {
		return nil, err
	}
	desc.Sentence = sentence
	nounOpt, has := nouns[sentence.Noun].(map[string]interface{})
	if !has {
		return nil, fmt.Errorf("Noun opt is missing or not a map.")
	}
	desc.nounOpt = nounOpt
	loc := speaker.VerbLocation(sentence.Noun, sentence.Verb)
	if loc == "" {
		return nil, fmt.Errorf("Verb %v in noun %v is not defined.", sentence.Verb, sentence.Noun)
	}
	desc.VerbLocation = loc
	return desc, nil
}

// Returns error if input is not valid according to the rules found in d.nouns
func (d *Descriptor) CreateInputs(filterCreator func(string, map[string]interface{})(iface.Filter,error)) ([]interface{}, map[string]interface{}, error) {
	module := mod.NewModule(d.VerbLocation)
	if !module.Exists() {
		return nil, nil, fmt.Errorf("Module named %v does not exist.", d.VerbLocation)
	}
	ins := module.Instance()
	verb := ins.Method(d.Sentence.Verb)
	an := verbinfo.NewAnalyzer(verb)
	ac := an.ArgCount()
	if len(d.Route.Queries) < ac {
		return nil, nil, fmt.Errorf("Not enough input to supply.")
	}
	if ac == 0 {
		return nil, nil, nil
	}
	fc := an.FilterCount()
	if fc > 0 && filterCreator == nil {
		return nil, nil, fmt.Errorf("filterCreator is needed but it is nil.")
	}
	var inp []interface{}
	var data map[string]interface{}
	source := []map[string]interface{}{}
	for _, v := range d.Route.Queries {
		source = append(source, v)
	}
	if an.NeedsData() {
		data = source[len(source)-1]
		if data == nil {
			data = map[string]interface{}{}		// !Important.
		}
	}
	if fc > 0 {
		if fc != 1 && len(source) != fc {
			return nil, nil, fmt.Errorf("Got %v inputs, but method %v needs only %v filters. Currently filters can only be reduced to 1.", len(source), d.Sentence.Verb, fc)
		}
		filters := []iface.Filter{}
		for i, v := range source {
			if d.Sentence.Verb != "Get" && d.Sentence.Verb != "GetSingle" && i == len(source)-1 {
				break
			}
			filt, err := filterCreator(d.Route.Words[i], v)
			if err != nil {
				return nil, nil, err
			}
			filters = append(filters, filt)
		}
		if len(filters) > 1 {
			filter := filters[0]
			red, err := filter.Reduce(filters[1:]...)
			if err != nil {
				return nil, nil, err
			}
			filters = []iface.Filter{red}
		}
		for _, v := range filters {
			inp = append(inp, v)
		}
	}
	return inp, data, nil
}