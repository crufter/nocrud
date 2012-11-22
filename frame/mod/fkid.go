package mod

import "github.com/opesun/nocrud/modules/fkid"

func init() {
	mods.register("fkid", fkid.C{})
}
