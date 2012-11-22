package lang_test

import (
	"github.com/opesun/nocrud/frame/lang"
	"github.com/opesun/nocrud/frame/misc/convert"
	"testing"
)

type Values map[string]interface{}

func (v Values) Add(i string, value interface{}) {
	convert.MapAdd(v, i, value)
}

func TestRoute(t *testing.T) {
	path := "/cars/comments"
	query := Values{}
	query.Add("make", "bmw")
	query.Add("engine", "4000")
	query.Add("1public", "true")
	route, err := lang.NewRoute(path, query)
	if err != nil {
		t.Fatal()
	}
	if len(route.Words) != 2 {
		t.Fatal()
	}
	if route.Words[0] != "cars" || route.Words[1] != "comments" {
		t.Fatal()
	}
	if route.Queries[0]["make"] == nil || route.Queries[0]["engine"] == nil || route.Queries[1]["public"] == nil {
		t.Fatal()
	}
}

func TestRoute1(t *testing.T) {
	path := "/x/y/z"
	query := Values{}
	route, err := lang.NewRoute(path, query)
	if err != nil {
		t.Fatal(err)
	}
	if len(route.Queries) != 3 {
		t.Fatal()
	}
}

func TestRoute2(t *testing.T) {
	path := "/x"
	query := Values{}
	query.Add("4hello", "this should fail")
	_, err := lang.NewRoute(path, query)
	if err == nil {
		t.Fatal()
	}
}

type MockSpeaker struct{}

func (m MockSpeaker) IsNoun(s string) bool {
	if s == "cars" || s == "comments" {
		return true
	}
	return false
}

func (m MockSpeaker) NounHasVerb(n, v string) bool {
	if n == "cars" && v == "Ignite" {
		return true
	}
	if n == "comments" && v == "Flame" {
		return true
	}
	return false
}

func TestSentence(t *testing.T) {
	path := "/cars/ignite"
	query := Values{}
	route, err := lang.NewRoute(path, query)
	if err != nil {
		t.Fatal()
	}
	speaker := MockSpeaker{}
	sentence, err := lang.NewSentence(route, speaker)
	if err != nil {
		t.Fatal(err)
	}
	if sentence.Noun != "cars" || sentence.Verb != "Ignite" || sentence.Redundant != "" {
		t.Fatal()
	}
}

func TestSentence1(t *testing.T) {
	path := "/cars/not-existing-verb"
	query := Values{}
	route, err := lang.NewRoute(path, query)
	if err != nil {
		t.Fatal()
	}
	speaker := MockSpeaker{}
	_, err = lang.NewSentence(route, speaker)
	if err == nil {
		t.Fatal()
	}
}

func TestSentence2(t *testing.T) {
	path := "/not-existing-noun/ignite"
	query := Values{}
	route, err := lang.NewRoute(path, query)
	if err != nil {
		t.Fatal()
	}
	speaker := MockSpeaker{}
	_, err = lang.NewSentence(route, speaker)
	if err == nil {
		t.Fatal()
	}
}

func TestSingle(t *testing.T) {
	path := "/cars/UHPHs2-Q6Q7Ey1gJ"
	query := Values{}
	route, err := lang.NewRoute(path, query)
	if err != nil {
		t.Fatal()
	}
	speaker := MockSpeaker{}
	sentence, err := lang.NewSentence(route, speaker)
	if err != nil {
		t.Fatal()
	}
	if sentence.Verb != "GetSingle" {
		t.Fatal()
	}
}

func TestURLEncoderUrlGet(t *testing.T) {
	path := "/cars"
	query := Values{}
	query.Add("favourites", "true")
	route, err := lang.NewRoute(path, query)
	if err != nil {
		t.Fatal()
	}
	speaker := MockSpeaker{}
	sentence, err := lang.NewSentence(route, speaker)
	if err != nil {
		t.Fatal()
	}
	if sentence.Verb != "Get" {
		t.Fatal()
	}
	urle := lang.NewURLEncoder(route, sentence)
	u1 := Values{}
	u1.Add("color", "red")
	u1.Add("color", "blue")
	u1.Add("quality", "very high")
	path, merged := urle.Url("paint", u1)
	if path != "cars/paint" {
		t.Fatal()
	}
	if merged["favourites"] != "true" || merged["1color"].([]interface{})[0] != "red" || merged["1color"].([]interface{})[1] != "blue" || merged["1quality"] != "very high" {
		t.Fatal()
	}
	if len(merged) != 3 {
		t.Fatal(merged)
	}
}

func TestURLEncoderUrlNonGet(t *testing.T) {
	path := "/cars/ignite"
	query := Values{}
	query.Add("favourites", "true")
	query.Add("1fake", "11")
	route, err := lang.NewRoute(path, query)
	if err != nil {
		t.Fatal()
	}
	speaker := MockSpeaker{}
	sentence, err := lang.NewSentence(route, speaker)
	if err != nil {
		t.Fatal()
	}
	if sentence.Verb != "Ignite" {
		t.Fatal()
	}
	urle := lang.NewURLEncoder(route, sentence)
	u1 := Values{}
	u1.Add("color", "red")
	u1.Add("color", "blue")
	u1.Add("quality", "very high")
	path, merged := urle.Url("paint", u1)
	if path != "cars/paint" {
		t.Fatal(path)
	}
	if merged["favourites"] != "true" || merged["1color"].([]interface{})[0] != "red" || merged["1color"].([]interface{})[1] != "blue" || merged["1quality"] != "very high" {
		t.Fatal(merged)
	}
}

func TestUrlEncoderForm(t *testing.T) {
	path := "/cars/UHPHs2-Q6Q7Ey1gJ/comments/flame"
	query := Values{}
	query.Add("favourites", "true")
	query.Add("1fake", "11")
	route, err := lang.NewRoute(path, query)
	if err != nil {
		t.Fatal()
	}
	speaker := MockSpeaker{}
	sentence, err := lang.NewSentence(route, speaker)
	if err != nil {
		t.Fatal()
	}
	if sentence.Verb != "Flame" && sentence.Noun != "comments" {
		t.Fatal()
	}
	urle := lang.NewURLEncoder(route, sentence)
	form := urle.Form("whatever-action")
	if form.KeyPrefix != "1" {
		t.Fatal(form.KeyPrefix)
	}
	if len(form.FilterFields) != 2 || form.FilterFields["favourites"] != "true" || form.FilterFields["1fake"] != "11" {
		t.Fatal()
	}
	if form.ActionPath != "cars/UHPHs2-Q6Q7Ey1gJ/comments/whatever-action" {
		t.Fatal()
	}
}

func TestUrlEncoderForm1(t *testing.T) {
	path := "/cars/ignite"
	query := Values{}
	route, err := lang.NewRoute(path, query)
	if err != nil {
		t.Fatal()
	}
	speaker := MockSpeaker{}
	sentence, err := lang.NewSentence(route, speaker)
	if err != nil {
		t.Fatal()
	}
	urle := lang.NewURLEncoder(route, sentence)
	form := urle.Form("whatever-action")
	if form.KeyPrefix != "1" {
		t.Fatal(form.KeyPrefix)
	}
}
