package main

import (
	"cli"
	"core"
	"fmt"
)

func NewTag(cmd cli.Cmd) {
	var tag core.Tag
	tag.Name = core.SureName(cmd.Pa[0])
	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "desc":
			tag.Desc = fl.V
		}
	}

	tag.ID = core.NewTag(tag.Name, tag.Desc)
	fmt.Printf("tag %d created\n", tag.ID)
}
func GetTag(cmd cli.Cmd) {
	ret := core.GetTag()
	fmt.Println(core.PrintTags(ret))
}
func GetTagByName(cmd cli.Cmd) {
	var id = core.SureID(core.GetTagByName(cmd.Pa[0]))
	fmt.Println(core.PrintTags([]int64{id}))
}
func GetTagByDesc(cmd cli.Cmd) {
	var ids = core.GetTagByDesc(cmd.Pa[0])
	fmt.Println(core.PrintTags(ids))
}
func AltTag(cmd cli.Cmd) {
	var desc *string = nil
	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "desc":
			desc = &fl.V
		}
	}

	core.SureID(core.AltTag(cmd.Pa[0], desc))
}
func AltTagRename(cmd cli.Cmd) {
	var old, new string
	old = cmd.Pa[0]
	new = cmd.Pa[1]

	core.SureName(new)

	core.AltRenameTag(old, new)
}
func DelTag(cmd cli.Cmd) {
	core.SureID(core.DelTag(cmd.Pa[0]))
}
