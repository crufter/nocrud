package session

import (
	iface "github.com/opesun/nocrud/frame/interfaces"
	"labix.org/v2/mgo"
)

type Session struct {
	session *mgo.Session
	opt     map[string]interface{}
	hooks   iface.Hooks
	newDb   func(session *mgo.Session, db *mgo.Database, opt map[string]interface{}, hooks iface.Hooks) iface.Db
}

// This is an ugly hack to get around the ugly circular dependency.
// We will get rid of that later.
func New(session *mgo.Session, opt map[string]interface{}, hooks iface.Hooks, newDb func(session *mgo.Session, db *mgo.Database, opt map[string]interface{}, hooks iface.Hooks) iface.Db) *Session {
	return &Session{
		session,
		opt,
		hooks,
		newDb,
	}
}

func (s *Session) Db(name string) (iface.Db, error) {
	// We should check if we have rights to select a given Db, and return error accordingly.
	// Currently, since the mgo session.Database involves no network communication, we will only receive errors when trying to do a DB operation.
	db := s.session.DB(name)
	return s.newDb(s.session, db, s.opt, s.hooks), nil
}
