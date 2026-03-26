package main

import (
	"cli"
	"core"
	"fmt"
)

func NewAcc(cmd cli.Cmd) {
	var acc core.Acc
	var err error

	acc.Name, err = core.SureName(cmd.Pa[0])

	if err != nil {
		panic(err)
	}

	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "desc":
			acc.Desc = fl.V
		}
	}

	acc.ID, err = core.NewAcc(acc.Name, acc.Desc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("account #%d created\n", acc.ID)
}

func GetAcc(cmd cli.Cmd) {
	ret := core.GetAcc()
	fmt.Println(core.PrintAccs(ret))
}
func GetAccByName(cmd cli.Cmd) {
	id, err := core.GetAccByName(cmd.Pa[0])
	if err != nil {
		panic(err)
	}

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

	core.AltAcc(cmd.Pa[0], desc)
}
func AltAccRename(cmd cli.Cmd) {
	var old, new string
	old = cmd.Pa[0]
	new = cmd.Pa[1]

	core.SureName(new)

	core.AltRenameAcc(old, new)
}
func DelAcc(cmd cli.Cmd) {
	core.DelAcc(cmd.Pa[0])
}
func GetAccChild(cmd cli.Cmd) {
	name := cmd.Pa[0]
	_, err := core.GetAccByName(cmd.Pa[0])
	if err != nil {
		panic(err)
	}

	ids := core.GetAccByPrefix(name)
	core.PrintAccs(ids)
}
