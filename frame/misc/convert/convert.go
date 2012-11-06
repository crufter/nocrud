package convert

import(
	"labix.org/v2/mgo/bson"
	"sort"
	"fmt"
)

// Cleans all bson.M s to map[string]interface{} s. Usually called on db query results.
// Will become obsolete when the mgo driver will return map[string]interface{} maps instead of bson.M ones.
func Clean(x interface{}) interface{} {
	if y, ok := x.(bson.M); ok {
		for key, val := range y {
			y[key] = Clean(val)
		}
		return (map[string]interface{})(y)
	} else if d, ok := x.(map[string]interface{}); ok {
		for key, val := range d {
			d[key] = Clean(val)
		}
		return d
	} else if z, ok := x.([]interface{}); ok {
		for i, v := range z {
			z[i] = Clean(v)
		}
	}
	return x
}

func Mapify(a map[string][]string) map[string]interface{} {
	ret := map[string]interface{}{}
	for i, v := range a {
		if len(v) == 1 {
			ret[i] = v[0]
		} else {
			sl := []interface{}{}
			for _, elem := range v {
				sl = append(sl, elem)
			}
			ret[i] = sl
		}
	}
	return ret
}

func createItem(key string, scheme interface{}, dat interface{}) map[string]interface{} {
	item := map[string]interface{}{"key": key}
	item["value"] = dat
	if sch, ok := scheme.(map[string]interface{}); ok {
		if typ, hast := sch["type"]; hast {
			item["type"] = typ
		}
		if disp, hasd := sch["disp"]; hasd {
			item["disp"] = disp
		}
	}
	return item
}

// Takes a dat map[string]interface{}, and puts every element of that which is defined in r to a slice, sorted by the keys ABC order.
// prior parameter can override the default abc ordering, so keys in prior will be the first ones in the slice, if those keys exist.
func abcKeys(scheme map[string]interface{}, dat map[string]interface{}, prior []string) []map[string]interface{} {
	ret := []map[string]interface{}{}
	already_in := map[string]struct{}{}
	for _, v := range prior {
		if _, contains := scheme[v]; contains {
			item := createItem(v, scheme[v], dat[v])
			ret = append(ret, item)
			already_in[v] = struct{}{}
		}
	}
	keys := []string{}
	for i, v := range scheme {
		// If the value is not false
		if boo, is_boo := v.(bool); !is_boo || boo == true {
			keys = append(keys, i)
		}
	}
	sort.Strings(keys)
	for _, v := range keys {
		if _, in := already_in[v]; !in {
			item := createItem(v, scheme[v], dat[v])
			ret = append(ret, item)
		}
	}
	return ret
}

// Takes an extraction/validation scheme, a document and from that creates a slice which can be easily displayed by a templating engine as a html form.
// Takes interface{}s and not map[string]interface{}s to include type checking here, and avoid that boilerplate in caller. 
func SchemeToFields(scheme interface{}, dat interface{}) ([]map[string]interface{}, error) {
	rm, rm_ok := scheme.(map[string]interface{})
	if !rm_ok {
		return nil, fmt.Errorf("Scheme is not a map[string]interface{}.")
	}
	datm, datm_ok := dat.(map[string]interface{})
	if !datm_ok && dat != nil {
		return nil, fmt.Errorf("Dat is not a map[string]interface{}.")
	}
	return abcKeys(rm, datm, []string{"title", "name", "slug"}), nil
}

// A more generic version of abcKeys. Takes a map[string]interface{} and puts every element of that into an []interface{}, ordered by keys alphabetically.
// TODO: find the intersecting parts between the two functions and refactor.
func OrderKeys(d map[string]interface{}) []interface{} {
	keys := []string{}
	for i, _ := range d {
		keys = append(keys, i)
	}
	sort.Strings(keys)
	ret := []interface{}{}
	for _, v := range keys {
		if ma, is_ma := d[v].(map[string]interface{}); is_ma {
			// RETHINK: What if a key field gets overwritten? Should we name it _key?
			ma["key"] = v
		}
		ret = append(ret, d[v])
	}
	return ret
}

// Converts all bson.ObjectId s to string. Usually called before displaying a database query result.
// Input is the result from the database.
func Recurs(v interface{}, converter func(interface{})(interface{},bool)) interface{} {
	switch value := v.(type) {
	case bson.M:
		for i, mem := range value {
			if conv, ok := converter(mem); ok {
				value[i] = conv
			} else {
				Recurs(mem, converter)
			}
		}
	case map[string]interface{}:
		for i, mem := range value {
			if conv, ok := converter(mem); ok {
				value[i] = conv
			} else {
				Recurs(mem, converter)
			}
		}
	case []map[string]interface{}:
		for i, mem := range value {
			if conv, ok := converter(mem); ok {
				value[i] = conv.(map[string]interface{})
			} else {
				Recurs(mem, converter)
			}
		}
	case []interface{}:
		for i, mem := range value {
			if conv, ok := converter(mem); ok {
				value[i] = conv
			} else {
				Recurs(mem, converter)
			}
		}
	}
	return v
}

func MapAdd(ma map[string]interface{}, key string, value interface{}) {
	val, has := ma[key]
	if !has {
		ma[key] = value
		return
	}
	slice, ok := val.([]interface{})
	if !ok {
		sl := []interface{}{}
		sl = append(sl, val)
		sl = append(sl, value)
		slice = sl
	} else {
		slice = append(slice, value)
	}
	ma[key] = slice
} 

func ListToMap(i ...interface{}) map[string]interface{} {
	r := map[string]interface{}{}
	for x:=0;x<len(i)-1; {
		MapAdd(r, i[x].(string), i[x+1])
		x = x+2
	}
	return r
}

func ToStringSlice(i ...interface{}) []string {
	ret := []string{}
	for _, v := range i {
		ret = append(ret, v.(string))
	}
	return ret
}