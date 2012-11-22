package document

import (
	iface "github.com/opesun/nocrud/frame/interfaces"
)

type Document struct {
	set  iface.Set
	data map[string]interface{}
	id   iface.Id
}

func New(set iface.Set, data map[string]interface{}) iface.Document {
	return &Document{
		set,
		data,
		data["_id"].(iface.Id),
	}
}

func (g *Document) Data() map[string]interface{} {
	return g.data
}

func (g *Document) Update(upd map[string]interface{}) error {
	q := map[string]interface{}{
		"_id": g.id,
	}
	return g.set.Update(q, upd)
}

func (g *Document) Remove() error {
	q := map[string]interface{}{
		"_id": g.id,
	}
	return g.set.Remove(q)
}

func (g *Document) Id() iface.Id {
	return g.id
}
