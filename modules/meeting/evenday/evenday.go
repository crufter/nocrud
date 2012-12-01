package evenday

import (
	"encoding/json"
	"fmt"
	"github.com/opesun/numcon"
	"strconv"
	"strings"
	"time"
)

type DayName int

const (
	OddMon DayName = iota
	OddTue
	OddWed
	OddThu
	OddFri
	OddSat
	OddSun
	EvenTue
	EvenWed
	EvenThu
	EvenFri
	EvenSat
	EvenSun
)

// n = seconds from Unix Epoch.
// The weeks starts on monday.
func DateToDayName(n int64) DayName {
	t := time.Unix(n, 0)
	_, week := t.ISOWeek()
	odd := week%2 == 1
	day := t.Weekday()
	// Correcting the "week starts on Sunday" approach of the time package.
	if day == time.Sunday {
		day = 6
	} else {
		day = day - 1
	}
	if odd {
		return DayName(day)
	}
	return DayName(day + 7)
}

func DateToMinute(n int64) Minute {
	hour :=	time.Unix(n, 0).Hour() - 1	// This was off by one, don't know why. 	
	min := 	time.Unix(n, 0).Minute()
	return Minute(int(hour)*60+int(min))
}

// Xth minute of a day. 0th minute is 0:00
type Minute int

func (m Minute) String() string {
	hour := m / 60
	minute := m % 60
	return strconv.Itoa(int(hour)) + ":" + strconv.Itoa(int(minute))
}

func ToMinute(hour, minute int) (Minute, error) {
	if hour < 0 || hour > 23 {
		return 0, fmt.Errorf("Hour is not proper.")
	}
	if minute < 0 || minute > 59 {
		return 0, fmt.Errorf("Minute is not proper.")
	}
	return Minute(hour*60 + minute), nil
}

// Converts strings like "8:10" or "08:10" to minutes like 8*60+10
func StringToMinute(s string) (Minute, bool) {
	spl := strings.Split(s, ":")
	if len(spl) != 2 {
		return 0, false
	}
	hour, err := strconv.Atoi(spl[0])
	if err != nil {
		return 0, false
	}
	minute, err := strconv.Atoi(spl[1])
	if err != nil {
		return 0, false
	}
	if hour > 23 || minute > 59 {
		return 0, false
	}
	return Minute(hour*60 + minute), true
}

// From and To fields were both unexported, but then the mgo driver can't serialize it...
type Interval struct {
	From Minute
	To   Minute
}

// Convenience function.
func ToInterval(fromHour, fromMinute, toHour, toMinute int) (*Interval, error) {
	fromMins, err := ToMinute(fromHour, fromMinute)
	if err != nil {
		return nil, err
	}
	toMins, err := ToMinute(toHour, toMinute)
	if err != nil {
		return nil, err
	}
	return NewInterval(fromMins, toMins)
}

// func (i *Interval) From() Minute {
// 	return i.From
// }
// 
// func (i *Interval) To() Minute {
// 	return i.To
// }

func (i *Interval) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"from": i.From,
		"to":   i.To,
	})
}

func (i *Interval) String() string {
	return i.From.String() + "-" + i.To.String()
}

// Returns true if interval a fits in interval b.
func InInterval(a, b *Interval) bool {
	return a.From >= b.From && a.To <= b.To
}

func NewInterval(from, to Minute) (*Interval, error) {
	if from > to {
		return nil, fmt.Errorf("From is greated than to.")
	}
	return &Interval{
		from,
		to,
	}, nil
}

// Handles both Minute and date inputs
func GenericToInterval(fromI, toI interface{}) (*Interval, error) {
	from, err := numcon.Int(fromI)
	if err != nil {
		return nil, err
	}
	to, err := numcon.Int(toI)
	if err != nil {
		return nil, err
	}
	if from > 1440 {
		from = int(DateToMinute(int64(from))) // Ouch...
	}
	if to > 1440 {
		to = int(DateToMinute(int64(to)))
	}
	interval, err := NewInterval(Minute(from), Minute(to))
	if err != nil {
		return nil, err
	}
	return interval, nil
}

func StringToInterval(s string) (*Interval, error) {
	fromTo := strings.Split(s, "-")
	if len(fromTo) != 2 {
		return nil, fmt.Errorf("Interval malformed.")
	}
	from, ok := StringToMinute(fromTo[0])
	if !ok {
		return nil, fmt.Errorf("From malformed.")
	}
	to, ok := StringToMinute(fromTo[1])
	if !ok {
		return nil, fmt.Errorf("To malformed.")
	}
	return &Interval{
		from,
		to,
	}, nil
}

// DaySchedule is a list of intervals when one is open to meetings.
type DaySchedule []*Interval

func (d DaySchedule) String() string {
	intStr := []string{}
	for _, v := range d {
		intStr = append(intStr, v.String())
	}
	return strings.Join(intStr, ", ")
}

func GenericToDaySchedule(a []interface{}) (DaySchedule, error) {
	ret := []*Interval{}
	for _, v := range a {
		m, ok := v.(map[string]interface{})
		if !ok {
			return DaySchedule{}, fmt.Errorf("Interval is not a map[string]interface{}.")
		}
		interval, err := GenericToInterval(m["from"], m["to"])
		if err != nil {
			return DaySchedule{}, err
		}
		ret = append(ret, interval)
	}
	return DaySchedule(ret), nil
}

// Converts a daystring "8:00-12:00, 13:00-15:00" to Intervals.
func StringToDaySchedule(s string) (DaySchedule, error) {
	spl := strings.Split(s, ",")
	ret := []*Interval{}
	for _, v := range spl {
		v = strings.Trim(v, " ")
		interval, err := StringToInterval(v)
		if err != nil {
			return DaySchedule{}, err
		}
		ret = append(ret, interval)
	}
	return DaySchedule(ret), nil
}

// Returns true if an interval fits into a Schedule.
func InSchedule(a *Interval, sch DaySchedule) bool {
	for _, v := range sch {
		if InInterval(a, v) {
			return true
		}
	}
	return false
}

type TimeTable [14]DaySchedule

func (tt *TimeTable) String() string {
	dayStr := []string{}
	for _, v := range tt {
		dayStr = append(dayStr, v.String())
	}
	return strings.Join(dayStr, ". ")
}

// Converts from JSON-like representation.
func GenericToTimeTable(g []interface{}) (*TimeTable, error) {
	ret := &TimeTable{}
	if len(g) != 14 {
		return nil, fmt.Errorf("Bad format.")
	}
	for i, v := range g {
		sl, ok := v.([]interface{})
		if !ok {
			return nil, fmt.Errorf("DaySchedule is not an []interface.")
		}
		ds, err := GenericToDaySchedule(sl)
		if err != nil {
			return nil, err
		}
		ret[i] = ds
	}
	return ret, nil
}

func StringsToTimeTable(s []string) (*TimeTable, error) {
	tt := &TimeTable{}
	for i, v := range s {
		v = strings.Trim(v, " \\n")
		ds, err := StringToDaySchedule(v)
		if err != nil {
			return nil, err
		}
		tt[i] = ds
	}
	return tt, nil
}

func StringToTimeTable(s string) (*TimeTable, error) {
	sl := strings.Split(s, ".")
	if len(s) != 14 {
		return nil, fmt.Errorf("Not a complete timetable.")
	}
	return StringsToTimeTable(sl)
}

func InTimeTable(dn DayName, i *Interval, tt *TimeTable) bool {
	return InSchedule(i, tt[dn])
}
