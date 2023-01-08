package main

import (
	"cli"
	"core"
	"fmt"
)

func NewAcc(cmd cli.Cmd) {
	var acc core.Acc
	acc.Name = core.SureName(cmd.Pa[0])
	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "desc":
			acc.Desc = fl.V
		}
	}

	acc.ID = core.NewAcc(acc.Name, acc.Desc)
	fmt.Printf("account #%d created\n", acc.ID)
}
func GetAcc(cmd cli.Cmd) {
	ret := core.GetAcc()
	fmt.Println(core.PrintAccs(ret))
}
func GetAccByName(cmd cli.Cmd) {
	var id = core.SureID(core.GetAccByName(cmd.Pa[0]))
	fmt.Println(core.PrintAccs([]int64{id}))
}
func GetAccByDesc(cmd cli.Cmd) {
	var ids = core.GetAccByDesc(cmd.Pa[0])
	fmt.Println(core.PrintAccs(ids))
}
func AltAcc(cmd cli.Cmd) {
	var desc *string = nil
	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "desc":
			desc = &fl.V
		}
	}

	core.SureID(core.AltAcc(cmd.Pa[0], desc))
}
func AltAccRename(cmd cli.Cmd) {
	var old, new string
	old = cmd.Pa[0]
	new = cmd.Pa[1]

	core.SureName(new)

	core.SureID(core.AltRenameAcc(old, new))
}
func DelAcc(cmd cli.Cmd) {
	core.SureID(core.DelAcc(cmd.Pa[0]))
}
func GetAccChild(cmd cli.Cmd) {
	name := cmd.Pa[0]
	core.SureID(core.GetAccByName(cmd.Pa[0]))

	ids := core.GetAccByPrefix(name)
	core.PrintAccs(ids)
}
