package mod

import "github.com/opesun/nocrud/modules/file_return_example"

func init() {
	mods.register("file_return_example", file_return_example.C{})
}
