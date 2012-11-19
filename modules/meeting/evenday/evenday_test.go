package evenday_test

import(
	"github.com/opesun/nocrud/modules/pretime/evenday"
	"testing"
	"strconv"
)

func TestStringToMinute(t *testing.T) {
	for hour:=0;hour<24;hour++{
		for min:=0;min<60;min++{
			hStr := strconv.Itoa(hour)
			mStr := strconv.Itoa(min)
			mins, ok := evenday.StringToMinute(hStr + ":" + mStr)
			if !ok {
				t.Fatal(hStr, mStr)
			}
			if mins != evenday.Minute(hour*60+min) {
				t.Fatal(mins, hour*60+min)
			}
		}
	}
}

func TestStringToMinuteLeadingZero(t *testing.T) {
	a := "8:10"
	b := "08:10"
	aM, ok := evenday.StringToMinute(a)
	if !ok {
		t.Fatal()
	}
	bM, ok := evenday.StringToMinute(b)
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

func toMinute(a, b int) evenday.Minute {
	val, err := evenday.ToMinute(a, b)
	if err != nil {
		panic(err)
	}
	return val
}

func TestStringToDaySchedule(t *testing.T) {
	ds, err := evenday.StringToDaySchedule("8:00-12:00, 13:30-17:40, 19:05-21:15")
	if err != nil {
		t.Fatal(err)
	}
	if ds[0].From() != toMinute(8, 0) {
		t.Fatal(ds[0])
	}
	if ds[0].To() != toMinute(12, 0) {
		t.Fatal(ds[0])
	}
	if ds[1].From() != toMinute(13, 30) {
		t.Fatal(ds[1])
	}
	if ds[1].To() != toMinute(17, 40) {
		t.Fatal(ds[1])
	}
	if ds[2].From() != toMinute(19, 5) {
		t.Fatal(ds[2])
	}
	if ds[2].To() != toMinute(21, 15) {
		t.Fatal(ds[2])
	}
}

func TestInInterval(t *testing.T) {
	a, err := evenday.ToInterval(8, 20, 9, 30)
	if err != nil {
		t.Fatal(err)
	}
	b, err := evenday.ToInterval(8, 30, 8, 50)
	if err != nil {
		t.Fatal(err)
	}
	c, err := evenday.ToInterval(11, 00, 12, 00)
	if err != nil {
		t.Fatal(err)
	}
	if !evenday.InInterval(b, a) {
		t.Fatal()
	}
	if evenday.InInterval(c, a) {
		t.Fatal()
	}
}

func TestInSchedule(t *testing.T) {
	ds, err := evenday.StringToDaySchedule("9:20-11:15, 13:20-15:30, 17:05-18:15")
	if err != nil {
		t.Fatal(err)
	}
	interval, err := evenday.ToInterval(9, 30, 10, 20)
	if err != nil {
		t.Fatal(err)
	}
	in := evenday.InSchedule(interval, ds)
	if !in {
		t.Fatal()
	}
	interval, err = evenday.ToInterval(8, 20, 9, 30)
	if err != nil {
		t.Fatal(err)
	}
	in = evenday.InSchedule(interval, ds)
	if in {
		t.Fatal()
	}
}

func toInterval(a, b, c, d int) *evenday.Interval {
	inte, err := evenday.ToInterval(a, b, c, d)
	if err != nil {
		panic(err)
	}
	return inte
}

func TestGenericToTimeTable(t *testing.T) {
	weekday := []interface{}{
		map[string]interface{}{
			"from": int(toMinute(8, 0)),
			"to": int(toMinute(17, 0)),
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
	tt, err := evenday.GenericToTimeTable(g)
	if err != nil {
		t.Fatal(err)
	}
	if !evenday.InTimeTable(evenday.OddMon, toInterval(8, 30, 9, 30), tt) {
		t.Fatal()
	}
	if evenday.InTimeTable(evenday.OddMon, toInterval(7, 30, 9, 30), tt) {
		t.Fatal()
	}
}