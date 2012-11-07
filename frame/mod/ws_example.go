package mod

import "github.com/opesun/nocrud/modules/ws_example"

func init() {
	mods.register("ws_example", ws_example.C{})
}