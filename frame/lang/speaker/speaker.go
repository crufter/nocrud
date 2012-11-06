package speaker

type Speaker struct {
	NounsToVerbs map[string]interface{}
	Has func(string, string) bool
	Fallback string
}

func New(a func(string, string) bool, b map[string]interface{}) *Speaker {
	return &Speaker{b, a, ""}
}

func (t *Speaker) IsNoun(a string) bool {
	_, has := t.NounsToVerbs[a]
	return has
}

func (t *Speaker) NounHasVerb(noun, verb string) bool {
	return t.VerbLocation(noun, verb) != ""
}

func (t *Speaker) verbLocation(noun, verb string) string {
	val, has := t.NounsToVerbs[noun].(map[string]interface{})
	if !has {
		return ""
	}
	for _, v := range val["composed_of"].([]interface{}) {
		if t.Has(v.(string), verb) {
			return v.(string)
		}
	}
	return ""
}

func (t *Speaker) VerbLocation(noun, verb string) string {
	loc := t.verbLocation(noun, verb)
	if loc == "" && t.Fallback != "" {
		loc = t.Fallback
	}
	return loc
}