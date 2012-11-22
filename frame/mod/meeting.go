package mod

import "github.com/opesun/nocrud/modules/meeting"

func init() {
	mods.register("meeting.entries", meeting.Entries{})
	mods.register("meeting.timeTable", meeting.TimeTable{})
}
