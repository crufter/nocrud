// Package user implements basic user functionality.
// - Registration, deletion, update, login, logout of users.
// - Building the user itself (if logged in), and putting it to uni.Dat["_user"].
package users

import (
	"fmt"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/modules/users/model"
)

type C struct {
	client iface.Client
	db     iface.Db
}

func (c *C) Init(ctx iface.Context) {
	c.client = ctx.Client()
	c.db = ctx.Db()
}

func (a *C) Insert(f iface.Filter, data map[string]interface{}) (iface.Id, error) {
	data["level"] = 100
	return user_model.RegisterUser(f, a.db, data)
}

func (a *C) InsertAdmin(f iface.Filter, data map[string]interface{}) (iface.Id, error) {
	err := hasAdmin(f)
	if err != nil {
		return nil, nil
	}
	data["level"] = 300
	return user_model.RegisterUser(f, a.db, data)
}

func (a *C) LoginForm() error {
	return nil
}

func (a *C) New() error {
	return nil
}

func hasAdmin(f iface.Filter) error {
	q := map[string]interface{}{
		"level": 300,
	}
	f.AddQuery(q)
	c, err := f.Count()
	if err != nil {
		return err
	}
	if c > 0 {
		return fmt.Errorf("Site already has an admin.")
	}
	return nil
}

func (a *C) NewAdmin(f iface.Filter) error {
	return hasAdmin(f)
}

func (a *C) Login(f iface.Filter, data map[string]interface{}) error {
	name := data["name"].(string)
	pass := data["password"].(string)
	// Maybe there could be a check here to not log in somebody who is already logged in.
	if user, err := user_model.FindLogin(f, name, pass); err == nil {
		a.client.StoreEncrypted("user", f.Subject()+"|"+user.Id().String())
	} else {
		return err
	}
	return nil
}

func (a *C) Logout() error {
	a.client.Unstore("user")
	return nil
}
