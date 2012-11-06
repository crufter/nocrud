package mod

import ca "github.com/opesun/nocrud/modules/custom_actions"

func init() {
	mods.register("custom_actions", ca.C{})
}