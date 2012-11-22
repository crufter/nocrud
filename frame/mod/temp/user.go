package mod

import "github.com/opesun/nocrud/modules/user"

func init() {
	mods.register("user", user.C{})
}
