package meeting

import(
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/modules/meeting/evenday"
	"fmt"
)

type C struct {
	user	iface.User
}

func (c *C) Init(ctx iface.Context) {
	c.user = ctx.User()
}

func toSS(sl []interface{}) []string {
	ret := []string{}
	for _, v := range sl {
		ret = append(ret, v.(string))
	}
	return ret
}

func (c *C) Entries() ([]map[string]interface{}, error) {
	return nil, nil
}

func (c *C) GetClosest() (evenday.Interval, error) {
	return evenday.Interval{}, nil
}

func (c *C) SaveEntry(data map[string]interface{}) error {
	return nil
}

func (c *C) ClientEntries(data map[string]interface{}) error {
	return nil
}

func (c *C) DeleteEntry(a iface.Filter) error {
	return nil
}

func (c *C) TimeTable() ([]string, error) {
	return nil, nil
}

func (c *C) SaveTimetable(a iface.Filter, data map[string]interface{}) error {
	if _, ok := c.user.Data()["professional"]; !ok {
		return fmt.Errorf("Only professionals can save timetables.")
	}
	//ttString, ok := data["timetable"].([]interface{})
	//if !ok {
	//	return fmt.Errorf("No timetable to save.")
	//}
	//err := validateTimetable(timetable)
	//if err != nil {
	//	return err
	//}
	count, err := a.Count()
	if err != nil {
		return err
	}
	m := map[string]interface{}{}
	m["created_by"] = c.user.Id()
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