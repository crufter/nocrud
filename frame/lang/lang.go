package lang

import(
	"strconv"
	"strings"
	"fmt"
	"regexp"
	"net/url"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/frame/misc/convert"
)

func ToCodeStyle(a string) string {
	a = strings.Replace(a, "-", " ", -1)
	a = strings.Replace(a, "_", " ", -1)
	a = strings.Title(a)
	return strings.Replace(a, " ", "", -1)
}

func ToURLStyle(a string) string {
	return back(a, "-")
}

func back(a, sep string) string {
	r := regexp.MustCompile("([A-Z])")
	res := r.ReplaceAll([]byte(a), []byte(" $1"))
	spl := strings.Split(string(res), " ")
	for i := range spl {
		spl[i] = strings.ToLower(spl[i])
	}
	return strings.Join(spl, sep)
}

func ToFileStyle(a string) string {
	return back(a, "_")
}

type Route struct {
	checked			int
	Words			[]string
	Queries			[]map[string]interface{}
}

type URLEncoder struct {
	r *Route
	s *Sentence
}

func NewURLEncoder(r *Route, s *Sentence) *URLEncoder {
	return &URLEncoder{r, s}
}

func loneId(a map[string]interface{}) bool {
	if len(a) != 1 {
		return false
	}
	if a["id"] == nil {
		return false
	}
	if slice, yepp := a["id"].([]interface{}); yepp {
		if len(slice) > 1 {
			return false
		}
	}
	return true
}

func extractLone(a map[string]interface{}) string {
	return a["id"].(string)
}

func (u *URLEncoder) actionPath(action_name string) string {
	var words []string
	rang := u.r.Words
	if u.s.Verb != "Get" && u.s.Verb != "GetSingle" {
		rang = rang[:len(rang)-1]
	}
	for i, v := range rang {
		words = append(words, v)
		if loneId(u.r.Queries[i]) {
			words = append(words, extractLone(u.r.Queries[i]))
		}
	}
	words = append(words, action_name)
	path := strings.Join(words, "/")
	return path
}

func mapToQueryString(a map[string]interface{}) string {
	u := url.Values{}
	for i, v := range a {
		switch t := v.(type) {
		case []interface{}:
			for _, v1 := range t {
				u.Add(i, fmt.Sprint(v1))
			}
		default:
			u.Add(i, fmt.Sprint(t))
		}
	}
	return u.Encode()
}

func (u *URLEncoder) UrlString(action_name string, input_params map[string]interface{}) string {
	path, merged := u.Url(action_name, input_params)
	qu := mapToQueryString(merged)
	if len(qu) > 0 {
		path = path+"?"+qu
	}
	return path
}

func (u *URLEncoder) Url(action_name string, input_params map[string]interface{}) (string, map[string]interface{}) {
	path := u.actionPath(action_name)
	var l []map[string]interface{}
	if u.s.Verb != "Get" && u.s.Verb != "GetSingle" {
		l = u.r.Queries
		l[len(l)-1] = input_params
	} else {
		l = append(u.r.Queries, input_params)
	}
	return path, EncodeQueries(l, true)
}

type Form struct {
	FilterFields	map[string]interface{}
	ActionPath		string
	KeyPrefix		string
}

func keyPrefix(q []map[string]interface{}) int {
	dec := 0
	for _, v := range q {
		if loneId(v) {
			dec++
		}
	}
	return len(q)-dec-1
}

func keyPrefixString(q []map[string]interface{}) string {
	return strconv.Itoa(keyPrefix(q))
}

func (u *URLEncoder) Form(action_name string) *Form {
	f := &Form{}
	f.ActionPath = u.actionPath(action_name)
	f.KeyPrefix = keyPrefixString(u.r.Queries)
	f.FilterFields = u.r.EncodeQueries(true)
	return f
}

func EncodeQueries(queries []map[string]interface{}, ignore_lone_id bool) map[string]interface{} {
	u := map[string]interface{}{}
	dec := 0
	for i, v := range queries {
		if ignore_lone_id && loneId(v) {
			dec++
			continue
		}
		var prefix string
		if i-dec != 0 {
			prefix = strconv.Itoa(i-dec)
		}
		for j, x := range v {
			if _, yepp := x.([]interface{}); yepp {
				for _, z := range x.([]interface{}) {
					convert.MapAdd(u, prefix+j, z)
				}
			} else {
				u[prefix+j] = x
			}
		}
	}
	return u
}

func (r *Route) EncodeQueries(ignore_lone_id bool) map[string]interface{} {
	return EncodeQueries(r.Queries, ignore_lone_id)
}

func (r *Route) Get() string {
	r.checked++
	return r.Words[len(r.Words)-r.checked]
}

func (r *Route) Got() int {
	return r.checked
}

func (r *Route) DropOne() {
	r.Words = r.Words[:len(r.Words)-1]
	r.Queries = r.Queries[:len(r.Queries)-1]
}

func sortParams(q map[string]interface{}) map[int]map[string]interface{} {
	sorted := map[int]map[string]interface{}{}
	for i, v := range q {
		num, err := strconv.Atoi(string(i[0]))
		nummed := false
		if err == nil {
			nummed = true
		} else {
			num = 0
		}
		if nummed {
			i = i[1:]
		}
		if _, has := sorted[num]; !has {
			sorted[num] = map[string]interface{}{}
		}
		if _, ok := v.([]interface{}); ok {
			for _, x := range v.([]interface{}) {
				convert.MapAdd(sorted[num], i, x)
			}
		} else {
			sorted[num][i] = v
		}
	}
	return sorted
}

func hasLargerThan(q map[int]map[string]interface{}, n int) bool {
	for i := range q {
		if i > n {
			return true
		}
	}
	return false
}

func nextIsId(current, next string) bool {
	//return next[1] == '-' && current[0] == next[0]
	return len(next) == 16
}

func extractId(next string) string {
	//return strings.Split(next, "-")[1]
	return next
}

// When creating a route, we essentially move input data from the path to the queries, eg.
// /cars/:id becomes /cars?id=:id
// and expanding the flattened query params, eg.
// /cars/comments?make=bmw&1date=today becomes /cars?make=bmw /comments?date=today
func NewRoute(path string, q map[string]interface{}) (*Route, error) {
	ps := strings.Split(path, "/")
	r := &Route{}
	r.Queries = []map[string]interface{}{}
	r.Words = []string{}
	if len(ps) < 1 {
		return r, fmt.Errorf("Wtf.")
	}
	ps = ps[1:]		// First one is empty string.
	sorted := sortParams(q)
	skipped := 0
	for i:=0;i<len(ps);i++ {
		v := ps[i]
		r.Words = append(r.Words, v)
		r.Queries = append(r.Queries, map[string]interface{}{})
		qi := len(r.Words)-1
		if len(ps) > i+1 {	// We are not at the end.
			next := ps[i+1]
			if nextIsId(v, next) {	// Id query in url., eg /users/fxARrttgFd34xdv7
				skipped++
				cq := r.Queries[qi]
				val, has := cq["id"]
				if !has {
					cq["id"] = extractId(next)
				} else {
					slice, ok := val.([]interface{})
					if !ok {
						sl := []interface{}{}
						sl = append(sl, val)
						sl = append(sl, extractId(next))
						slice = sl
					} else {
						slice = append(slice, extractId(next))
					}
					r.Queries[qi]["id"] = slice
				}
				i++
				continue
			}
		}
		r.Queries[qi] = sorted[qi-skipped]
	}
	if hasLargerThan(sorted, len(r.Words)-1) {
		return nil, fmt.Errorf("Unnecessary sorted params.")
	}
	return r, nil
}

type Sentence struct{
	Noun, Verb, Redundant string
}

// The construction of a sentence consists of analyzing parts of the route with the help of a speaker,
// who can recognize nouns, and verbs related to that noun.
//
// Rethink later: there can be a certain ambivalence in the way the nouns (subjects) determine the
// location for a verb.
// Eg (retard example):
// /users/posts/whatevers/delete-everything		=>		The verb "DeleteEverything" is not related only to "whatevers", but rather it is
// a standalone method acting on arguments...
func NewSentence(r *Route, speaker iface.Speaker) (*Sentence, error) {
	s := &Sentence{}
	if len(r.Words) == 1 {
		s.Noun = r.Words[0]
		if loneId(r.Queries[0]) {
			s.Verb = "GetSingle"
		} else {
			s.Verb = "Get"
		}
		return s, nil
	}
	unstable := r.Get()
	must_be_noun := r.Get()
	l := len(r.Words)
	if speaker.IsNoun(unstable) {
		if loneId(r.Queries[l-1]) {
			s.Verb = "GetSingle"
		} else {
			s.Verb = "Get"
		}
		s.Noun = unstable
	} else if speaker.NounHasVerb(must_be_noun, ToCodeStyle(unstable)) {
		s.Verb = ToCodeStyle(unstable)
		s.Noun = must_be_noun
	} else {
		s.Redundant = unstable
		// A noun is singular if it has exactly one query param, the id.
		if !loneId(r.Queries[l-2]) {
			return nil, fmt.Errorf("Plural nouns can't have redundant information.")
		}
		r.DropOne()
		s.Verb = "GetSingle"
		s.Noun = must_be_noun
	}
	if !speaker.IsNoun(s.Noun) {
		return nil, fmt.Errorf("%v is not a noun.", s.Noun)
	}
	return s, nil
}