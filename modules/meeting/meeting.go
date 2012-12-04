package meeting

import (
	"fmt"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/modules/meeting/evenday"
)

func isProfessional(u iface.User) bool {
	_, ok := u.Data()["professional"]
	return ok
}

type shared struct {
	db						iface.Db
	userId 					iface.Id
	userIsProfessional		bool
	optDoc					iface.NestedData
	timeTableColl			string
	intervalColl			string
	gotOptions				bool
}

type Entries struct {
	shared
}

func (e *Entries) Init(ctx iface.Context) {
	e.db = ctx.Db()
	e.userId = ctx.User().Id()
	e.userIsProfessional = isProfessional(ctx.User())
	e.optDoc = ctx.Options().Document()
	e.timeTableColl = "timeTables"
	e.intervalColl = "intervals"
}

func (e *Entries) getOptions(resource string) {
	if e.gotOptions {
		return
	}
	e.gotOptions = true
	ttColl, ok := e.optDoc.GetStr("nouns." + resource + ".options.timeTableColl")
	if ok {
		e.timeTableColl = ttColl
	}
	iColl, ok := e.optDoc.GetStr("nouns." + resource + ".options.intervalColl")
	if ok {
		e.intervalColl = iColl
	}
}

// Returns the closest interval on the same day to the given interval.
func (e *Entries) GetClosest(a iface.Filter, data map[string]interface{}) (evenday.Interval, error) {
	e.getOptions(a.Subject())
	//prof, err := e.db.ToId(data["professional"].(string))
	//if err != nil {
	//	return evenday.Interval{}, err
	//}
	//from := data["from"].(int64)
	//length := data["length"].(int64)
	//err := e.intervalIsValid(data, length)
	//if err != nil {
	//	return err
	//}
	//to := from + length * 60
	//day := evenday.DateToDayname
	return evenday.Interval{}, nil
}

func (e *Entries) getTimeTable(prof iface.Id) (*evenday.TimeTable, error) {
	ttFilter, err := e.db.NewFilter(e.timeTableColl, nil)
	if err != nil {
		return nil, err
	}
	timeTableQ := map[string]interface{}{
		"createdBy": prof,
	}
	ttFilter.AddQuery(timeTableQ)
	ttC, err := ttFilter.Count()
	if err != nil {
		return nil, err
	}
	if ttC != 1 {
		return nil, fmt.Errorf("There are multiple timetables.")
	}
	timeTables, err := ttFilter.Find()
	if err != nil {
		return nil, err
	}
	timeTable, err := evenday.GenericToTimeTable(timeTables[0]["timeTable"].([]interface{}))
	if err != nil {
		return nil, err
	}
	return timeTable, nil
}

// Checks if the timeTable is ok and the interval fits into the timeTable.
func (e *Entries) okAccordingToTimeTable(data map[string]interface{}, from, to int64) error {
	dayN := evenday.DateToDayName(from)
	prof, err := e.db.ToId(data["professional"].(string))
	if err != nil {
		return err
	}
	interval, err := evenday.GenericToInterval(from, to)
	if err != nil {
		return err
	}
	timeTable, err := e.getTimeTable(prof)
	if err != nil {
		return err
	}
	if !evenday.InTimeTable(dayN, interval, timeTable) {
		return fmt.Errorf("Interval does not fit into timeTable.")
	}
	return nil
}

func (e *Entries) othersAlreadyTook(a iface.Filter, from, to int64) error {
	entryQ := map[string]interface{}{
		"$or": []interface{}{
			map[string]interface{}{
				"from": map[string]interface{}{
					"$gt": from,
					"$lt": to,
				},
			},
			map[string]interface{}{
				"to": map[string]interface{}{
					"$gt": from,
					"$lt": to,
				},
			},
		},
	}
	a.AddQuery(entryQ)
	eC, err := a.Count()
	if err != nil {
		return err
	}
	if eC > 0 {
		return fmt.Errorf("That time is already taken.")
	}
	return nil
}

func (e *Entries) intervalIsValid(data map[string]interface{}, length int64) error {
	prof, err := e.db.ToId(data["professional"].(string))
	if err != nil {
		return err
	}
	q := map[string]interface{}{
		"professional": prof,
		"length": length,
	}
	iFilter, err := e.db.NewFilter("intervals", q)
	if err != nil {
		return err
	}
	c, err := iFilter.Count()
	if err != nil {
		return err
	}
	if c != 1 {
		return fmt.Errorf("Interval %v is not defined.", length)
	}
	return nil
}

func (e *Entries) Insert(a iface.Filter, data map[string]interface{}) error {
	e.getOptions(a.Subject())
	from := data["from"].(int64)
	length := data["length"].(int64)
	err := e.intervalIsValid(data, length)
	if err != nil {
		return err
	}
	to := from + length * 60
	err = e.okAccordingToTimeTable(data, from, to)
	if err != nil {
		return err
	}
	err = e.othersAlreadyTook(a, from, to)
	if err != nil {
		return err
	}
	i := map[string]interface{}{
		"createdBy": e.userId,
		"from": data["from"],
		"to": to,
	}
	return a.Insert(i)
}

func (e *Entries) Delete(a iface.Filter) error {
	return nil
}

// With the current design, a professional can't view his own entries posted to account of an other professional.
func (e *Entries) ProcessQuery(q map[string]interface{}) {
	if e.userIsProfessional {
		q["professional"] = e.userId
	} else {
		q["createdBy"] = e.userId
	}
}

func (e *Entries) Install(o iface.Document, resource string) error {
	upd := map[string]interface{}{
		"$addToSet": map[string]interface{}{
			"Hooks." + resource + "ProcessQuery": []interface{}{"entries", "ProcessQuery"},
		},
	}
	return o.Update(upd)
}

func (e *Entries) Uninstall(o iface.Document, resource string) error {
	upd := map[string]interface{}{
		"$pull": map[string]interface{}{
			"Hooks." + resource + "ProcessQuery": []interface{}{"entries", "ProcessQuery"},
		},
	}
	return o.Update(upd)
}

type TimeTable struct {
	shared
}

func (tt *TimeTable) Init(ctx iface.Context) {
	tt.db = ctx.Db()
	tt.userId = ctx.User().Id()
	tt.userIsProfessional = isProfessional(ctx.User())
	tt.optDoc = ctx.Options().Document()
	tt.timeTableColl = "timeTables"
	tt.intervalColl = "intervals"
}

func toSS(sl []interface{}) []string {
	ret := []string{}
	for _, v := range sl {
		ret = append(ret, v.(string))
	}
	return ret
}

func (tt *TimeTable) Save(a iface.Filter, data map[string]interface{}) error {
	if !tt.userIsProfessional {
		return fmt.Errorf("Only professionals can save timetables.")
	}
	ssl := toSS(data["timeTable"].([]interface{}))
	timeTable, err := evenday.StringsToTimeTable(ssl)
	if err != nil {
		return err
	}
	count, err := a.Count()
	if err != nil {
		return err
	}
	m := map[string]interface{}{}
	m["createdBy"] = tt.userId
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
	q["createdBy"] = tt.userId
}

func (tt *TimeTable) Install(o iface.Document, resource string) error {
	upd := map[string]interface{}{
		"$addToSet": map[string]interface{}{
			"Hooks." + resource + "ProcessQuery": []interface{}{"timetable", "ProcessQuery"},
		},
	}
	return o.Update(upd)
}

func (tt *TimeTable) Uninstall(o iface.Document, resource string) error {
	upd := map[string]interface{}{
		"$pull": map[string]interface{}{
			"Hooks." + resource + "ProcessQuery": []interface{}{"timetable", "ProcessQuery"},
		},
	}
	return o.Update(upd)
}
