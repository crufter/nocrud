package basics

import (
	iface "github.com/opesun/nocrud/frame/interfaces"
)

type Basics struct {
	Hooks iface.Hooks
	Db    iface.Db
}

type QueryInfo struct {
	Count   int
	Skipped int
	Limited int
	Sorted  []string
}

func (b *Basics) Get(a iface.Filter) ([]map[string]interface{}, *QueryInfo, error) {
	list, err := a.Find()
	if err != nil {
		return nil, nil, err
	}
	count, err := a.Count()
	if err != nil {
		return nil, nil, err
	}
	return list, &QueryInfo{
		count, a.Modifiers().Skip(),
		a.Modifiers().Limit(),
		a.Modifiers().Sort(),
	}, nil
}

func (b *Basics) GetSingle(a iface.Filter) (map[string]interface{}, error) {
	return a.FindOne()
}

func (b *Basics) Insert(a iface.Filter, data map[string]interface{}) (iface.Id, error) {
	id := b.Db.NewId()
	data["_id"] = id
	err := a.Insert(data)
	if err != nil {
		return nil, err
	}
	if b.Hooks != nil {
		q := map[string]interface{}{
			"_id": id,
		}
		filt := a.Clone().AddQuery(q)
		b.Hooks.Select("Inserted").Fire(filt)
		b.Hooks.Select(a.Subject() + "Inserted").Fire(filt)
	}
	return id, nil
}

func (b *Basics) Update(a iface.Filter, data map[string]interface{}) error {
	upd := map[string]interface{}{
		"$set": data,
	}
	err := a.Update(upd)
	if err != nil {
		return err
	}
	if b.Hooks != nil {
		b.Hooks.Select("Updated").Fire(a)
		b.Hooks.Select(a.Subject() + "Updated").Fire(a)
	}
	return nil
}

func (b *Basics) UpdateAll(a iface.Filter, data map[string]interface{}) error {
	upd := map[string]interface{}{
		"$set": data,
	}
	_, err := a.UpdateAll(upd)
	return err
}

func (b *Basics) Remove(a iface.Filter) error {
	return a.Remove()
}

func (b *Basics) RemoveAll(a iface.Filter) error {
	_, err := a.RemoveAll()
	return err
}
