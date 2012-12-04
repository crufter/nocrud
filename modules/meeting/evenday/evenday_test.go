package evenday

import (
	"strconv"
	"testing"
	"time"
)

func TestStringToMinute(t *testing.T) {
	for hour := 0; hour < 24; hour++ {
		for min := 0; min < 60; min++ {
			hStr := strconv.Itoa(hour)
			mStr := strconv.Itoa(min)
			mins, ok := StringToMinute(hStr + ":" + mStr)
			if !ok {
				t.Fatal(hStr, mStr)
			}
			if mins != Minute(hour*60+min) {
				t.Fatal(mins, hour*60+min)
			}
		}
	}
}

func TestStringToMinuteLeadingZero(t *testing.T) {
	a := "8:10"
	b := "08:10"
	aM, ok := StringToMinute(a)
	if !ok {
		t.Fatal()
	}
	bM, ok := StringToMinute(b)
	if !ok {
		t.Fatal()
	}
	if aM == 0 {
		t.Fatal()
	}
	if bM == 0 {
		t.Fatal()
	}
	if aM != bM {
		t.Fatal()
	}
}

func toMinute(a, b int) Minute {
	val, err := ToMinute(a, b)
	if err != nil {
		panic(err)
	}
	return val
}

func testStringToDaySchedule(t *testing.T, dd dsData, iter int) {
	ds, err := StringToDaySchedule(dd.ser)
	if err != nil {
		t.Fatal(iter, err)
	}
	if len(ds.intervals) != len(dd.checks) {
		t.Fatal(iter, "Bad test data")
	}
	for i, v := range ds.intervals {
		cks := dd.checks[i]
		fromMin := toMinute(cks[0], cks[1])
		if v.From != fromMin {
			t.Fatal(iter, "Bad from, got:", fromMin, "want:", v.From)
		}
		toMin := toMinute(cks[2], cks[3])
		if v.To != toMin {
			t.Fatal(iter, "Bad to, got:", toMin, "want:", v.To)
		}
	}
}

type dsData struct {
	ser		string
	checks	[][4]int
}

func TestStringToDaySchedule(t *testing.T) {
	inp := []dsData {
		{ "8:00-12:00, 13:30-17:40, 19:05-21:15", [][4]int{ [4]int{8, 0, 12, 0}, [4]int{13, 30, 17, 40}, [4]int{19, 05, 21, 15} }},
	}
	for i, v := range inp {
		testStringToDaySchedule(t, v, i)
	}
}

func ti(a, b, c, d int) *Interval {
	ret, err := ToInterval(a, b, c, d)
	if err != nil {
		panic(err)
	}
	return ret
}

func testInInterval(t *testing.T, dat inIntervalData, iter int) {
	a := ti(dat.a[0], dat.a[1], dat.a[2], dat.a[3])
	b := ti(dat.b[0], dat.b[1], dat.b[2], dat.b[3])
	if InInterval(a, b) != dat.check {
		t.Fatal(iter)
	}
}

type inIntervalData struct {
	a		[4]int
	b		[4]int
	check	bool
}

func TestInInterval(t *testing.T) {
	inp := []inIntervalData {
		{ [4]int{8, 20, 9, 30}, [4]int{8, 20, 9, 30}, true },
		{ [4]int{11, 00, 12, 00}, [4]int{8, 20, 9, 30}, false },
		{ [4]int{8, 30, 8, 40}, [4]int{8, 20, 9, 30}, true },
	}
	for i, v := range inp {
		testInInterval(t, v, i)
	}
}

func testTouchesInterval(t *testing.T, dat touchesIntervalData, iter int) {
	a := ti(dat.a[0], dat.a[1], dat.a[2], dat.a[3])
	b := ti(dat.b[0], dat.b[1], dat.b[2], dat.b[3])
	if TouchesInterval(a, b) != dat.check {
		t.Fatal(iter)
	}
}

type touchesIntervalData struct {
	a		[4]int
	b		[4]int
	check 	bool
}

func TestTouchesInterval(t *testing.T) {
	inp := []touchesIntervalData{
		{ [4]int{8, 11, 8, 56}, [4]int{8, 29, 8, 55}, true },
		{ [4]int{8, 20, 8, 50}, [4]int{8, 50, 9, 30}, false },
		{ [4]int{8, 20, 8, 51}, [4]int{8, 50, 9, 30}, true },
	}
	for i, v := range inp {
		testTouchesInterval(t, v, i)
	}
}

func TestInDaySchedule(t *testing.T) {
	ds, err := StringToDaySchedule("9:20-11:15, 13:20-15:30, 17:05-18:15")
	if err != nil {
		t.Fatal(err)
	}
	interval, err := ToInterval(9, 20, 11, 15)
	if err != nil {
		t.Fatal(err)
	}
	in := InDaySchedule(interval, ds)
	if !in {
		t.Fatal()
	}
	interval, err = ToInterval(8, 20, 9, 30)
	if err != nil {
		t.Fatal(err)
	}
	in = InDaySchedule(interval, ds)
	if in {
		t.Fatal()
	}
}

func TestInDaySchedule1(t *testing.T) {
	ds, err := StringToDaySchedule("9:20-11:15, 13:20-15:30, 17:05-18:15")
	if err != nil {
		t.Fatal(err)
	}
	interval, err := ToInterval(9, 30, 9, 50)
	if err != nil {
		t.Fatal(err)
	}
	if !InDaySchedule(interval, ds) {
		t.Fatal()
	}
}

func toInterval(a, b, c, d int) *Interval {
	inte, err := ToInterval(a, b, c, d)
	if err != nil {
		panic(err)
	}
	return inte
}

func TestGenericToTimeTable(t *testing.T) {
	weekday := []interface{}{
		map[string]interface{}{
			"from": int(toMinute(8, 0)),
			"to":   int(toMinute(17, 0)),
		},
	}
	weekend := []interface{}{}
	g := []interface{}{
		weekday,
		weekday,
		weekday,
		weekday,
		weekday,
		weekend,
		weekend,
		weekday,
		weekday,
		weekday,
		weekday,
		weekday,
		weekend,
		weekend,
	}
	tt, err := GenericToTimeTable(g)
	if err != nil {
		t.Fatal(err)
	}
	if !InTimeTable(OddMon, toInterval(8, 30, 9, 30), tt) {
		t.Fatal()
	}
	if InTimeTable(OddMon, toInterval(7, 30, 9, 30), tt) {
		t.Fatal()
	}
}

func TestDayToDayName(t *testing.T) {
	d, err := time.Parse("2006-01-02 15:04", "2012-11-30 10:52")
	if err != nil {
		t.Fatal(err)
	}
	dayName := DateToDayName(d.Unix())
	if dayName != 11 {
		t.Fatal(dayName)
	}
}

func TestDateToMinute(t *testing.T) {
	d, err := time.Parse("2006-01-02 15:04", "2012-11-30 10:52")
	if err != nil {
		t.Fatal(err)
	}
	minute := DateToMinute(d.Unix())
	if minute != 10*60+52 {
		t.Fatal(minute)
	}
}

func testAdvise(t *testing.T, a adviseData, iter int) {
	open, err := StringToDaySchedule(a.open)
	if err != nil {
		t.Fatal(iter, err)
	}
	taken, err := StringToDaySchedule(a.taken)
	if err != nil {
		t.Fatal(iter, err)
	}
	adv := NewAdvisor(open, taken)
	if adv.minuteSteps != 5 {
		t.Fatal(iter, "This test assumes 5 minutes stepping of the advisor.")
	}
	adv.BackwardsToo = a.backwards
	adv.Amount(a.amt)
	interv, err := StringToInterval(a.adviseFor)
	if err != nil {
		t.Fatal(iter, err)
	}
	res := adv.Advise(interv)
	if len(res) != len(a.checks) {
		t.Fatal(iter, len(res), len(a.checks))
	}
	for i, v := range a.checks {
		c, err := StringToInterval(v)
		if err != nil {
			t.Fatal(err)
		}
		if !res[i].Eq(c) {
			t.Fatal(iter, "Want: ", c, " got: ", res[i])
		}
	}
}

type adviseData struct {
	open		string
	taken		string
	adviseFor	string
	backwards	bool
	amt			int
	checks		[]string
}

func TestAdvise(t *testing.T) {
	all := []adviseData {
		{ "9:20-11:15, 13:20-15:30, 17:05-18:15", "9:30-10:20, 13:55-14:10", "9:30-9:50", false, 1, []string{ "10:20-10:40"}},
		{ "8:00-12:00", "8:10-8:40", "8:20-8:30", false, 1, []string{ "8:40-8:50" }},
		{ "8:00-12:00", "8:05-8:20, 8:29-8:55, 9:05-9:55", "8:01-8:46", false, 1, []string{ "9:56-10:41" }}, 
	}
	for i, v := range all {
		testAdvise(t, v, i)
	}
}