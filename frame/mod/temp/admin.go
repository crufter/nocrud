package mod

import ad "github.com/opesun/nocrud/modules/admin"

func init() {
	mods.register("admin", ad.C{})
}
