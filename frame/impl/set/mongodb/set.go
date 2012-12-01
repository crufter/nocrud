package set

import (
	"encoding/base64"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/frame/misc/convert"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func New(db *mgo.Database, coll string) iface.Set {
	return &Set{
		db,
		coll,
		0,
		0,
		nil,
	}
}

type Set struct {
	db    *mgo.Database
	coll  string
	skip  int
	limit int
	sort  []string
}

func (s *Set) Skip(i int) {
	s.skip = i
}

func (s *Set) Limit(i int) {
	s.limit = i
}

func (s *Set) Sort(str ...string) {
	s.sort = str
}

func (s *Set) FindOne(q map[string]interface{}) (map[string]interface{}, error) {
	q = convert.Recurs(q, fromGeneral).(map[string]interface{})
	var res interface{}
	err := s.db.C(s.coll).Find(q).One(&res)
	if err != nil {
		return nil, err
	}
	ret := convert.Clean(res).(map[string]interface{})
	return convert.Recurs(ret, toGeneral).(map[string]interface{}), nil
}

func (s *Set) Count(q map[string]interface{}) (int, error) {
	q = convert.Recurs(q, fromGeneral).(map[string]interface{})
	return s.db.C(s.coll).Find(q).Count()
}

func (s *Set) Find(q map[string]interface{}) ([]map[string]interface{}, error) {
	q = convert.Recurs(q, fromGeneral).(map[string]interface{})
	c := s.db.C(s.coll).Find(q)
	if s.skip != 0 {
		c.Skip(s.skip)
	}
	if s.limit != 0 {
		c.Limit(s.limit)
	}
	if len(s.sort) > 0 {
		c.Sort(s.sort...)
	}
	var res []interface{}
	err := c.All(&res)
	if err != nil {
		return nil, err
	}
	isl := convert.Clean(res).([]interface{})
	ret := []map[string]interface{}{}
	for _, v := range isl {
		ret = append(ret, v.(map[string]interface{}))
	}
	ret = convert.Recurs(ret, toGeneral).([]map[string]interface{})
	return ret, nil
}

func (s *Set) Insert(d map[string]interface{}) error {
	d = convert.Recurs(d, fromGeneral).(map[string]interface{})
	return s.db.C(s.coll).Insert(d)
}

func (s *Set) Update(q map[string]interface{}, updQuery map[string]interface{}) error {
	q = convert.Recurs(q, fromGeneral).(map[string]interface{})
	updQuery = convert.Recurs(updQuery, fromGeneral).(map[string]interface{})
	return s.db.C(s.coll).Update(q, updQuery)
}

func (s *Set) UpdateAll(q map[string]interface{}, updQuery map[string]interface{}) (int, error) {
	q = convert.Recurs(q, fromGeneral).(map[string]interface{})
	updQuery = convert.Recurs(updQuery, fromGeneral).(map[string]interface{})
	chi, err := s.db.C(s.coll).UpdateAll(q, updQuery)
	return chi.Updated, err
}

func (s *Set) Remove(q map[string]interface{}) error {
	q = convert.Recurs(q, fromGeneral).(map[string]interface{})
	return s.db.C(s.coll).Remove(q)
}

func (s *Set) RemoveAll(q map[string]interface{}) (int, error) {
	q = convert.Recurs(q, fromGeneral).(map[string]interface{})
	chi, err := s.db.C(s.coll).RemoveAll(q)
	return chi.Removed, err
}

func fromGeneral(a interface{}) (interface{}, bool) {
	val, ok := a.(iface.Id)
	if !ok {
		return a, false
	}
	return decodeIdP(val.String()), true
}

func toGeneral(a interface{}) (interface{}, bool) {
	val, ok := a.(bson.ObjectId)
	if !ok {
		return a, false
	}
	return &Id{
		val,
	}, true
}

func (s *Set) Name() string {
	return s.coll
}

type Id struct {
	i bson.ObjectId
}

func (id *Id) String() string {
	return base64.URLEncoding.EncodeToString([]byte(id.i))
}

func (id *Id) IAmAnId() {
}

func NewId() iface.Id {
	return &Id{
		bson.NewObjectId(),
	}
}

func ToId(encodedForm string) iface.Id {
	return &Id{
		decodeIdP(encodedForm),
	}
}

// This can be problematic, easily triggers false positives.
func IsId(encodedForm string) bool {
	if len(encodedForm) == 16 {
		return true
	}
	return false
}

func decodeId(s string) (bson.ObjectId, error) {
	val, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		panic("Can't decode id: " + s + " " + err.Error())
	}
	return bson.ObjectId(val), nil
}

func decodeIdP(s string) bson.ObjectId {
	val, err := decodeId(s)
	if err != nil {
		panic(err)
	}
	return val
}
