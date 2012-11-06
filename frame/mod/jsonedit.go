package mod

import "github.com/opesun/nocrud/modules/jsonedit"

func init() {
	mods.register("jsonedit", jsonedit.C{})
}