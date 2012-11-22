package db

import (
	"fmt"
	"github.com/opesun/jsonp"
	"github.com/opesun/nocrud/frame/impl/filter"
	"github.com/opesun/nocrud/frame/impl/session/mongodb"
	"github.com/opesun/nocrud/frame/impl/set/mongodb"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"labix.org/v2/mgo"
)

type Db struct {
	session *mgo.Session
	db      *mgo.Database
	opt     map[string]interface{}
	hooks   iface.Hooks
}

func New(session *mgo.Session, db *mgo.Database, opt map[string]interface{}, hooks iface.Hooks) iface.Db { // Returns iface.Db and not *Db to break the circular dependency problem.
	return &Db{
		session,
		db,
		opt,
		hooks,
	}
}

func (d *Db) NewFilter(c string, m map[string]interface{}) (iface.Filter, error) {
	s := set.New(d.db, c)
	filterScheme, _ := jsonp.GetM(d.opt, "nouns."+c+".filterScheme")
	return filter.New(s, d.hooks, filterScheme, m)
}

func (d *Db) ToId(i string) (iface.Id, error) {
	return set.ToId(i), nil
}

func (d *Db) NewId() iface.Id {
	return set.NewId()
}

func (d *Db) Session() (iface.Session, error) {
	if d.session == nil {
		return nil, fmt.Errorf("Seems like you don't have access to the session.")
	}
	return session.New(d.session, d.opt, d.hooks, New), nil
}
