package mod

import "github.com/opesun/nocrud/modules/terminal"

func init() {
	mods.register("terminal", terminal.C{})
}