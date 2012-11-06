package mod

import "github.com/opesun/nocrud/modules/file"

func init() {
	mods.register("file", file.C{})
}