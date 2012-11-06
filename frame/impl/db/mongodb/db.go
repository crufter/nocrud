package db

import(
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/frame/impl/filter"
	"github.com/opesun/nocrud/frame/impl/set/mongodb"
	"labix.org/v2/mgo"
)

type Db struct {
	db			*mgo.Database
	opt			map[string]interface{}
	hooks		iface.Hooks
}

func New(db *mgo.Database, opt map[string]interface{}, hooks iface.Hooks) *Db {
	return &Db{
		db,
		opt,
		hooks,
	}
}

func (d *Db) NewFilter(c string, m map[string]interface{}) (iface.Filter, error) {
	s := set.New(d.db, c)
	return filter.New(s, d.hooks, m)
}

func (d *Db) ToId(i string) (iface.Id, error) {
	return set.ToId(i), nil
}

func (d *Db) NewId() iface.Id {
	return set.NewId()
}