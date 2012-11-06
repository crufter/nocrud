package mod

import "github.com/opesun/nocrud/modules/fulltext"

func init() {
	mods.register("fulltext", fulltext.C{})
}