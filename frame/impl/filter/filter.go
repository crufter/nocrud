package filter

import (
	"fmt"
	"github.com/opesun/nocrud/frame/impl/document"
	"github.com/opesun/nocrud/frame/impl/set/mongodb" // ...
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/sanitize"
)

type Mods struct {
	skip  int
	limit int
	sort  []string
}

func (m *Mods) Skip() int {
	return m.skip
}

func (m *Mods) Limit() int {
	return m.limit
}

func (m *Mods) Sort() []string {
	return m.sort
}

// Parents are separated from the query, because they are not used only at querying (FindOne, Find, Update, UpdateAll, Remove, RemoveAll),
// but at Insert too.
type Filter struct {
	set         iface.Set
	mods        *Mods
	parentField map[string]string     // collection_name => fieldname, being used at reduction, shortcoming: 1 collection to 1 fieldname only...
	parents     map[string][]iface.Id // fieldnames => parent ids
	query       map[string]interface{}
	hooks       iface.Hooks
	scheme      map[string]interface{}
}

func NewSimple(set iface.Set, hooks iface.Hooks, scheme map[string]interface{}) (*Filter, error) {
	return &Filter{
		set:     set,
		parents: map[string][]iface.Id{},
		hooks:   hooks,
		scheme:  scheme,
	}, nil
}

func New(set iface.Set, hooks iface.Hooks, scheme, input map[string]interface{}) (*Filter, error) {
	d := processQuery(hooks, set.Name(), scheme, input)
	f := &Filter{
		set:  set,
		mods: d.mods,
		//parentField:	d.parentField,
		query:   d.query,
		parents: map[string][]iface.Id{},
		hooks:   hooks,
		scheme:  scheme,
	}
	return f, nil
}

func (f *Filter) Visualize() {
	fmt.Println("<<<")
	fmt.Println("fmod", f.mods)
	fmt.Println("parents", f.parents)
	fmt.Println("query", f.query)
	fmt.Println(">>>")
}

func (f *Filter) Reduce(a ...iface.Filter) (iface.Filter, error) {
	l := len(a)
	if l == 0 {
		return &Filter{}, fmt.Errorf("Nothing to reduce.")
	}
	var prev iface.Filter
	prev = f
	for _, v := range a {
		ids, err := prev.Ids()
		if err != nil {
			return &Filter{}, err
		}
		v.AddParents("_"+prev.Subject(), ids)
		prev = v
	}
	return prev, nil
}

// Information coming from url.Values/map
type data struct {
	query       map[string]interface{}
	mods        *Mods
	parentField string
}

// Special fields in query:
// parentf, sort, limit, skip, page
func processQuery(hooks iface.Hooks, coll string, scheme, inp map[string]interface{}) *data {
	d := &data{}
	if inp == nil {
		inp = map[string]interface{}{}
	}
	intSch := map[string]interface{}{
		"type": "int",
	}
	sch := map[string]interface{}{
		//"parentf": 1,
		"sort":  1,
		"skip":  intSch,
		"limit": intSch,
		"page":  intSch,
	}
	ex, err := sanitize.New(sch)
	if err != nil {
		panic(err)
	}
	dat, err := ex.Extract(inp)
	if err != nil {
		panic(err)
	}
	for i := range sch {
		delete(inp, i)
	}
	mods := &Mods{}
	//if dat["parentf"] != nil {
	//	d.parentField = dat["parentf"].(string)
	//}
	if dat["skip"] != nil {
		mods.skip = int(dat["skip"].(int64))
	}
	if dat["limit"] != nil {
		mods.limit = int(dat["limit"].(int64))
	} else {
		mods.limit = 20
	}
	if dat["page"] != nil {
		page := int(dat["page"].(int64))
		mods.skip = (page - 1) * mods.limit
	}
	if dat["sort"] != nil {
		mods.sort = []string{ dat["sort"].(string) }
	}
	d.mods = mods
	if hooks != nil {
		hooks.Select("ProcessQuery").Fire(inp) // We should let the subscriber now the subject name maybe.
		hooks.Select(coll + "ProcessQuery").Fire(inp)
	}
	if scheme != nil && len(scheme) != 0 {
		ex, err = sanitize.New(scheme)
		if err != nil {
			panic(err)
		}
		hooks.Select("SanitizerMangler").Fire(ex)
		dat, err := ex.Extract(inp)
		if err != nil {
			panic(err)
		}
		inp = dat
	}
	d.query = toQuery(inp)
	return d
}

func _append(vi []interface{}, i *string, x interface{}) []interface{} {
	if *i == "id" {
		*i = "_id"
		vi = append(vi, set.ToId(x.(string)))
	} else {
		vi = append(vi, x)
	}
	return vi
}

// map => mongodb query map
func toQuery(a map[string]interface{}) map[string]interface{} {
	r := map[string]interface{}{}
	for i, v := range a {
		if i[0] == '$' {
			r[i] = v
			continue
		}
		var vi []interface{}
		if slice, ok := v.([]interface{}); ok {
			for _, x := range slice {
				vi = _append(vi, &i, x)
			}
		} else {
			vi = _append(vi, &i, v)
		}
		if len(vi) > 1 { // Ex: {"$and": [{"fulltext": ^"whateverr"}, {...}]}
			r[i] = map[string]interface{}{
				"$in": vi,
			}
		} else {
			r[i] = vi[0]
		}
	}
	return r
}

func (f *Filter) Clone() iface.Filter {
	newM := map[string]interface{}{}
	for i, v := range f.query {
		newM[i] = v
	}
	return &Filter{
		set:			f.set,
		mods:			&*f.mods,
		parentField:	f.parentField,
		query:			newM,
		parents:		f.parents,
	}
}

func (f *Filter) Modifiers() iface.Modifiers {
	return f.mods
}

func (f *Filter) AddQuery(q map[string]interface{}) iface.Filter {
	query := processQuery(f.hooks, f.set.Name(), f.scheme, q).query
	for i, v := range f.query {
		query[i] = v
	}
	f.query = query
	return f
}

func mergeQuery(q map[string]interface{}, p map[string][]iface.Id) map[string]interface{} {
	r := map[string]interface{}{}
	for i, v := range q {
		r[i] = v
	}
	for i, v := range p {
		r[i] = map[string]interface{}{
			"$in": v,
		}
	}
	return r
}

func mergeInsert(ins map[string]interface{}, p map[string][]iface.Id) map[string]interface{} {
	r := map[string]interface{}{}
	for i, v := range ins {
		r[i] = v
	}
	for i, v := range p {
		r[i] = v
	}
	return r
}

func (f *Filter) FindOne() (map[string]interface{}, error) {
	q := mergeQuery(f.query, f.parents)
	return f.set.FindOne(q)
}

func (f *Filter) Find() ([]map[string]interface{}, error) {
	if f.mods.skip != 0 {
		f.set.Skip(f.mods.skip)
	}
	if f.mods.limit != 0 {
		f.set.Limit(f.mods.limit)
	}
	if len(f.mods.sort) > 0 {
		f.set.Sort(f.mods.sort...)
	}
	q := mergeQuery(f.query, f.parents)
	return f.set.Find(q)
}

func (f *Filter) SelectOne() (iface.Document, error) {
	q := mergeQuery(f.query, f.parents)
	data, err := f.set.FindOne(q)
	if err != nil {
		return nil, err
	}
	return document.New(f.set, data), nil
}

func (f *Filter) Iterate(callback func(iface.Document) error) error {
	f.set.Limit(0)
	q := mergeQuery(f.query, f.parents)
	dataz, err := f.set.Find(q)
	if err != nil {
		return err
	}
	for _, data := range dataz {
		doc := document.New(f.set, data)
		err := callback(doc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Filter) Insert(d map[string]interface{}) error {
	i := mergeInsert(d, f.parents)
	return f.set.Insert(i)
}

func (f *Filter) Update(upd_query map[string]interface{}) error {
	q := mergeQuery(f.query, f.parents)
	return f.set.Update(q, upd_query)
}

func (f *Filter) UpdateAll(upd_query map[string]interface{}) (int, error) {
	q := mergeQuery(f.query, f.parents)
	return f.set.UpdateAll(q, upd_query)
}

func (f *Filter) Subject() string {
	return f.set.Name()
}

func (f *Filter) Count() (int, error) {
	q := mergeQuery(f.query, f.parents)
	return f.set.Count(q)
}

func (f *Filter) AddParents(fieldname string, a []iface.Id) {
	slice, ok := f.parents[fieldname]
	if !ok {
		f.parents[fieldname] = []iface.Id{}
		slice = []iface.Id{}
	}
	slice = append(slice, a...)
	f.parents[fieldname] = slice
}

func (f *Filter) Ids() ([]iface.Id, error) {
	if val, has := f.query["id"]; has && len(f.query) == 1 && len(f.parents) == 1 {
		ids := val.(map[string]interface{})["$in"].([]iface.Id)
		return ids, nil
	}
	q := mergeQuery(f.query, f.parents)
	f.set.Limit(0)
	docs, err := f.set.Find(q)
	if err != nil {
		return nil, err
	}
	ret := []iface.Id{}
	for _, v := range docs {
		ret = append(ret, v["_id"].(iface.Id))
	}
	return ret, nil
}

func (f *Filter) Remove() error {
	q := mergeQuery(f.query, f.parents)
	return f.set.Remove(q)
}

func (f *Filter) RemoveAll() (int, error) {
	q := mergeQuery(f.query, f.parents)
	return f.set.RemoveAll(q)
}
