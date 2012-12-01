package meeting

import (
	iface "github.com/opesun/nocrud/frame/interfaces"
	mongoDb "github.com/opesun/nocrud/frame/impl/db/mongodb"
	emptyHooks "github.com/opesun/nocrud/frame/mocks/hooks"
	"github.com/opesun/nocrud/frame/impl/nesteddata"
	"labix.org/v2/mgo"
	"testing"
	"time"
)

var session *mgo.Session
var testDb *mgo.Database

var db iface.Db

var e *Entries
var tt *TimeTable

var newE func() iface.Filter
var newTT func() iface.Filter
var newI func() iface.Filter

var proId iface.Id

func init() {
	s, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	d := s.DB("test")
	err = d.DropDatabase()
	if err != nil {
		panic(err)
	}
	session = s
	testDb = d
	db = mongoDb.New(session, testDb, map[string]interface{}{}, emptyHooks.New())
	emptyOptDoc := nesteddata.New(map[string]interface{}{})
	clientId := db.NewId()
	e = &Entries{
		shared{
			db,
			clientId,
			false,
			emptyOptDoc,
			"timeTables",
			"intervals",
			false,
		},
	}
	proId = db.NewId()
	tt = &TimeTable{
		shared{
			db,
			proId,
			true,
			emptyOptDoc,
			"timeTables",
			"intervals",
			false,
		},
	}
	newFilter := func(coll string) iface.Filter {
		f, err := db.NewFilter(coll, nil)
		if err != nil {
			panic(err)
		}
		return f
	}
	newE = func() iface.Filter {
		return newFilter("entries")
	}
	newTT = func() iface.Filter {
		return newFilter("timeTables")
	}
	newI = func() iface.Filter {
		return newFilter("intervals")
	}
}

func TestEmptyTimeTables(t *testing.T) {
	clientInp := map[string]interface{}{
		"from": int64(1),
		"length": int64(20),
		"professional": proId.String(),
	}
	err := e.Insert(newE(), clientInp)
	if err == nil {
		t.Fatal("Without timetables, no insert should be possible.")
	}
}

func TestTimeTableSave(t *testing.T) {
	proInp := map[string]interface{}{
		"timeTable": []interface{}{
			"8:00-12:00, 13:00-17:00",
			"8:00-12:00, 13:00-17:00",
			"8:00-12:00, 13:00-17:00",
			"8:00-12:00, 13:00-17:00",
			"8:00-12:00, 13:00-17:00",
			"8:00-12:00",
			"0:00-0:00",
			"8:00-12:00, 13:00-17:00",
			"8:00-12:00, 13:00-17:00",
			"8:00-12:00, 13:00-17:00",
			"8:00-12:00, 13:00-17:00",
			"8:00-12:00, 13:00-17:00",
			"8:00-12:00",
			"0:00-0:00",
		},
	}
	err := tt.Save(newTT(), proInp)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBasicEntryInsert(t *testing.T) {
	fromT, err := time.Parse("2006-01-02 15:04", "2012-11-30 14:20")
	if err != nil {
		panic(err)
	}
	clientInp1 := map[string]interface{}{
		"from": fromT.Unix(),
		"length": int64(30),
		"professional": proId.String(),
	}
	err = e.Insert(newE(), clientInp1)
	if err == nil {
		t.Fatal("This test should fail because no interval is defined yet.")
	}
	i := map[string]interface{}{
		"professional": proId,
		"length": 30,
	}
	err = newI().Insert(i)
	if err != nil {
		t.Fatal(err)
	}
	err = e.Insert(newE(), clientInp1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestConflicting(t *testing.T) {
	fromT1, err := time.Parse("2006-01-02 15:04", "2012-11-30 14:49")
	if err != nil {
		t.Fatal(err)
	}
	clientInp2 := map[string]interface{}{
		"from": fromT1.Unix(),
		"length": int64(30),
		"professional": proId.String(),
	}
	c, err := newE().Count()
	if err != nil {
		t.Fatal(err)
	}
	if c != 1 {
		t.Fatal(c)
	}
	err = e.Insert(newE(), clientInp2)
	if err == nil {
		t.Fatal("This should conflict.")
	}
}

func TestUpperBoundary(t *testing.T) {
	fromT1, err := time.Parse("2006-01-02 15:04", "2012-11-30 14:50")
	if err != nil {
		t.Fatal(err)
	}
	clientInp2 := map[string]interface{}{
		"from": fromT1.Unix(),
		"length": int64(30),
		"professional": proId.String(),
	}
	c, err := newE().Count()
	if err != nil {
		t.Fatal(err)
	}
	if c != 1 {
		t.Fatal(c)
	}
	err = e.Insert(newE(), clientInp2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLowerBoundary(t *testing.T) {
	fromT1, err := time.Parse("2006-01-02 15:04", "2012-11-30 13:50")
	if err != nil {
		t.Fatal(err)
	}
	clientInp2 := map[string]interface{}{
		"from": fromT1.Unix(),
		"length": int64(30),
		"professional": proId.String(),
	}
	c, err := newE().Count()
	if err != nil {
		t.Fatal(err)
	}
	if c != 2 {
		t.Fatal(c)
	}
	err = e.Insert(newE(), clientInp2)
	if err != nil {
		t.Fatal(err)
	}
}