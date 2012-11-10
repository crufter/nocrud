package fulltext

import(
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/slugify"
	"labix.org/v2/mgo/bson"
	"strings"
	"fmt"
)

type C struct {
}

// TODO: rethink. TODO: Add number support.
// Walks an entire JSON tree recursively, and converts everything to string it can find.
func toStringRecursively(i interface{}) []string {
	switch val := i.(type) {
	case map[string]interface{}:
		ret := []string{}
		for _, v := range val {
			ret = append(ret, toStringRecursively(v)...)
		}
		return ret
	case []interface{}:
		ret := []string{}
		for _, v := range val {
			ret = append(ret, toStringRecursively(v)...)
		}
		return ret
	case string:
		return []string{val}
	}
	return []string{}
}

func filterDupes(s []string) []string {
	ret := []string{}
	c := map[string]struct{}{}
	for _, v := range s {
		if _, has := c[v]; !has {
			c[v] = struct{}{}
			ret = append(ret, v)
		}
	}
	return ret
}

func filterTooShort(s []string, min_len int) []string {
	ret := []string{}
	for _, v := range s {
		if len(v) >= min_len {
			ret = append(ret, v)
		}
	}
	return ret
}

func simpleFulltext(non_split []string) []string {
	split := []string{}
	for _, v := range non_split {
		split = append(split, strings.Split(v, " ")...)
	}
	slugified := []string{}
	for _, v := range split {
		slugified = append(slugified, strings.Trim(slugify.S(v), ",.:;"))
	}
	slugified = filterDupes(slugified)
	return filterTooShort(slugified, 3)
}

func (c *C) updateFromDoc(doc map[string]interface{}) map[string]interface{} {
	non_split := toStringRecursively(doc)
	fullt := simpleFulltext(non_split)
	upd := map[string]interface{}{
		"$set": map[string]interface{}{
			"fulltext": fullt,
		},
	}
	return upd
}

func (c *C) SaveFulltext(a iface.Filter) error {
	doc, err := a.FindOne()
	if err != nil {
		return err
	}
	upd := c.updateFromDoc(doc)
	return a.Update(upd)
}

func (c *C) RegenerateFulltext(a iface.Filter) error {
	cb := func(g iface.Document) error {
		upd := c.updateFromDoc(g.Data())
		return g.Update(upd)
	}
	return a.Iterate(cb)
}

func generateKeywords(s string) []string {
	split := strings.Split(s, " ")
	slugified := []string{}
	for _, v := range split {
		slugified = append(slugified, strings.Trim(slugify.S(v), ",.:;"))
	}
	return slugified
}

// Generates [{"fulltext": \^keyword1\}, {"fulltext": \^keyword2\}]
// With this query we can create a good enough full text search, which can search at the beginning of the keywords.
// We could write regexes which searches in the middle of the words too, but that query could not uzilize the btree indexes of mongodb.
// This solution must be efficient, assuming mongodb does the expected sane things: utilizing indexes with ^ regexes, "$and" queries and arrays.
func GenerateQuery(s string) []interface{} {
	sl := generateKeywords(s)
	and := []interface{}{}
	for _, v := range sl {
		and = append(and, map[string]interface{}{
			"fulltext": bson.RegEx{
				Pattern: "^" + v,
			},
		})
	}
	return and
}

func (c *C) ProcessMap(inp map[string]interface{}) {
	val, has := inp["search"].(string)
	if !has {
		return
	}
	inp["$and"] = GenerateQuery(val)
	delete(inp, "search")
}

func (c *C) Install(o iface.Document, s string) error {
	upd := map[string]interface{}{
		"$addToSet": map[string]interface{}{
			fmt.Sprintf("Hooks.%vInserted", s): []interface{}{"fulltext", "SaveFulltext"},
			fmt.Sprintf("Hooks.%vUpdated", s): []interface{}{"fulltext", "SaveFulltext"},
			"Hooks.ProcessMap": "fulltext",
		},
	}
	return o.Update(upd)
}

func (c *C) Uninstall(o iface.Document, s string) error {
	upd := map[string]interface{}{
		"$pull": map[string]interface{}{
			fmt.Sprintf("Hooks.%vInserted", s): []interface{}{"fulltext", "SaveFulltext"},
			fmt.Sprintf("Hooks.%vUpdated", s): []interface{}{"fulltext", "SaveFulltext"},
			"Hooks.ProcessMap": "fulltext",
		},
	}
	return o.Update(upd)
}