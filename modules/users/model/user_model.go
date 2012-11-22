package user_model

import (
	"crypto/sha1"
	"fmt"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"io"
)

func FindLogin(a iface.Filter, name, password string) (iface.Document, error) {
	encoded_pass := hashPass(password)
	q := map[string]interface{}{
		"name":     name,
		"password": encoded_pass,
	}
	return a.AddQuery(q).SelectOne()
}

func hashPass(pass string) string {
	h := sha1.New()
	io.WriteString(h, pass)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func RegisterUser(a iface.Filter, db iface.Db, user map[string]interface{}) (iface.Id, error) {
	user["password"] = hashPass(user["password"].(string))
	if _, has := user["level"]; !has {
		user["level"] = 100
	}
	user_id := db.NewId()
	user["_id"] = user_id
	err := a.Insert(user)
	if err != nil {
		return nil, fmt.Errorf("Name is not unique.")
	}
	return user_id, nil
}
