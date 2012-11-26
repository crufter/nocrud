package meeting

import (
	"fmt"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/modules/meeting/evenday"
)

type Entries struct {
	user iface.User
}

func (e *Entries) Init(ctx iface.Context) {
	e.user = ctx.User()
}

func (e *Entries) GetClosest() (evenday.Interval, error) {
	return evenday.Interval{}, nil
}

func (e *Entries) SaveEntry(data map[string]interface{}) error {
	return nil
}

func (e *Entries) ClientEntries(data map[string]interface{}) error {
	return nil
}

func (e *Entries) DeleteEntry(a iface.Filter) error {
	return nil
}

func (e *Entries) ProcessQuery(q map[string]interface{}) {
	q["createdBy"] = e.user.Id()
}

func (e *Entries) Install(o iface.Document, resource string) error {
	upd := map[string]interface{}{
		"$addToSet": map[string]interface{}{
			"Hooks." + resource + "ProcessQuery": []interface{"entries", "ProcessQuery"},
		},
	}
	return o.Update(upd)
}

func (e *Entries) Uninstall(o iface.Document, resource string) error {
	upd := map[string]interface{}{
		"$pull": map[string]interface{}{
			"Hooks." + resource + "ProcessQuery": []interface{"entries", "ProcessQuery"},
		},
	}
	return o.Update(upd)
}

type TimeTable struct {
	user iface.User
}

func (tt *TimeTable) Init(ctx iface.Context) {
	tt.user = ctx.User()
}

func toSS(sl []interface{}) []string {
	ret := []string{}
	for _, v := range sl {
		ret = append(ret, v.(string))
	}
	return ret
}

func (tt *TimeTable) Save(a iface.Filter, data map[string]interface{}) error {
	if _, ok := tt.user.Data()["professional"]; !ok {
		return fmt.Errorf("Only professionals can save timetables.")
	}
	ssl := toSS(data["timetable"].([]interface{}))
	timeTable, err := evenday.StringsToTimeTable(ssl)
	if err != nil {
		return err
	}
	count, err := a.Count()
	if err != nil {
		return err
	}
	m := map[string]interface{}{}
	m["createdBy"] = tt.user.Id()
	m["timeTable"] = timeTable
	if count == 0 {
		return a.Insert(m)
	} else if count == 1 {
		doc, err := a.SelectOne()
		if err != nil {
			return err
		}
		return doc.Update(m)
	}
	return fmt.Errorf("Too many timetables in the database.")
}

func (tt *TimeTable) ProcessQuery(q map[string]interface{}) {
	q["createdBy"] = tt.user.Id()
}

func (tt *TimeTable) Install(o iface.Document, resource string) error {
	upd := map[string]interface{}{
		"$addToSet": map[string]interface{}{
			"Hooks." + resource + "ProcessQuery": []interface{"timetable", "ProcessQuery"},
		},
	}
	return o.Update(upd)
}

func (tt *TimeTable) Uninstall(o iface.Document, resource string) error {
	upd := map[string]interface{}{
		"$pull": map[string]interface{}{
			"Hooks." + resource + "ProcessQuery": []interface{"timetable", "ProcessQuery"},
		},
	}
	return o.Update(upd)
}
