package mod

import te "github.com/opesun/nocrud/modules/template_editor"

func init() {
	mods.register("template_editor", te.C{})
}