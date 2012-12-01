package user

import (
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/numcon"
	"strings"
)

type User struct {
	iface.Document
	level     int
	languages []string
}

func New(db iface.Db, client iface.Client) *User {
	user, err := _new(db, client)
	if err != nil {
		return emptyUser(client)
	}
	return user
}

func _new(db iface.Db, client iface.Client) (*User, error) {
	uidI, err := client.GetDecrypted("user")
	if err != nil {
		return nil, err
	}
	uidStr := uidI.(string)
	coll := "users"
	spl := strings.Split(uidStr, "|")
	if len(spl) > 1 {
		coll = spl[0]
		uidStr = spl[1]
	}
	f, err := db.NewFilter(coll, nil)
	if err != nil {
		return nil, err
	}
	id, err := db.ToId(uidStr)
	if err != nil {
		return nil, err
	}
	q := map[string]interface{}{
		"_id": id,
	}
	f.AddQuery(q)
	userDoc, err := f.SelectOne()
	if err != nil {
		return nil, err
	}
	var langs []string
	if langz, ok := userDoc.Data()["languages"].([]interface{}); ok {
		for _, v := range langz {
			langs = append(langs, v.(string))
		}
	} else if client.Languages() != nil {
		langs = client.Languages()
	} else {
		langs = []string{"en"}
	}
	return &User{
		userDoc,
		numcon.IntP(userDoc.Data()["level"]),
		langs,
	}, nil
}

func (u *User) Languages() []string {
	return u.languages
}

func (u *User) Level() int {
	return u.level
}

// When no user cookie is found, or there was a problem during building the user,
// we proceed with an empty user.
func emptyUser(client iface.Client) *User {
	langs := client.Languages()
	if langs == nil {
		langs = []string{"en"}
	}
	return &User{
		nil,
		0,
		langs,
	}
}
