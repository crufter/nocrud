package speaker_test

import(
	"testing"
	"github.com/opesun/nocrud/frame/lang/speaker"
)

func moduleHasVerb(n, v string) bool {
	if n == "testModuleA" && v == "verbA" {
		return true
	}
	if n == "testModuleB" && v == "verbB" {
		return true
	}
	return false
}

func Test0 (t *testing.T) {
	mockNounOpt := map[string]interface{}{
		"cars": map[string]interface{}{
			"composed_of": []interface{}{"testModuleA"},
		},
	}
	spkr := speaker.New(moduleHasVerb, mockNounOpt)
	if !spkr.IsNoun("cars") {
		t.Fatal()
	}
	if !spkr.NounHasVerb("cars", "verbA") {
		t.Fatal()
	}
}

func Test1 (t *testing.T) {
	mockNounOpt := map[string]interface{}{
		"cars": map[string]interface{}{
			"composed_of": []interface{}{"testModuleA"},
		},
	}
	spkr := speaker.New(moduleHasVerb, mockNounOpt)
	if spkr.IsNoun("comments") {
		t.Fatal()
	}
	if spkr.NounHasVerb("cars", "verbB") {
		t.Fatal()
	}
}

func TestVerbLocation (t *testing.T) {
	mockNounOpt := map[string]interface{}{
		"cars": map[string]interface{}{
			"composed_of": []interface{}{"testModuleA"},
		},
	}
	spkr := speaker.New(moduleHasVerb, mockNounOpt)
	if spkr.VerbLocation("cars", "verbA") != "testModuleA" {
		t.Fatal()
	}
	if spkr.VerbLocation("x", "y") != "" {
		t.Fatal()
	}
	spkr.Fallback = "fakeModule"
	if spkr.VerbLocation("x", "y") != "fakeModule" {
		t.Fatal()
	}
}


