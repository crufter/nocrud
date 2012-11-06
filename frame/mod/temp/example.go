package mod

import "github.com/opesun/nocrud/modules/example"

func init() {
	mods.register("example", example.C{})
}