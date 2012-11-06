package mod

import "github.com/opesun/nocrud/modules/skeleton"

func init() {
	mods.register("skeleton", skeleton.C{})
}