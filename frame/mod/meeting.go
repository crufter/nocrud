package mod

import "github.com/opesun/nocrud/modules/meeting"

func init() {
	mods.register("meeting", meeting.C{})
}