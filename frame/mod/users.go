package mod

import "github.com/opesun/nocrud/modules/users"

func init() {
	mods.register("users", users.C{})
}
