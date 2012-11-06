package mod

import c "github.com/opesun/nocrud/modules/content"

func init() {
	mods.register("content", c.C{})
}