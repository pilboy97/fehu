package main

import (
	"cli"
	"core"
	"fmt"
)

func NewAcc(cmd cli.Cmd) error {
	var acc core.Acc
	var err error

	acc.Name, err = core.SureName(cmd.Pa[0])
	if err != nil {
		return err
	}

	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "desc":
			acc.Desc = fl.V
		}
	}

	acc.ID, err = core.NewAcc(acc.Name, acc.Desc)
	if err != nil {
		return err
	}
	fmt.Printf("account #%d created\n", acc.ID)
	return nil
}

func GetAcc(cmd cli.Cmd) error {
	ret, err := core.GetAcc()
	if err != nil {
		return err
	}
	fmt.Println(core.PrintAccs(ret))
	return nil
}

func GetAccByName(cmd cli.Cmd) error {
	id, err := core.GetAccByName(cmd.Pa[0])
	if err != nil {
		return err
	}
	fmt.Println(core.PrintAccs([]int64{id}))
	return nil
}

func GetAccByDesc(cmd cli.Cmd) error {
	ids, err := core.GetAccByDesc(cmd.Pa[0])
	if err != nil {
		return err
	}
	fmt.Println(core.PrintAccs(ids))
	return nil
}

func GetAccChild(cmd cli.Cmd) error {
	name := cmd.Pa[0]
	if _, err := core.GetAccByName(name); err != nil {
		return err
	}

	ids, err := core.GetAccByPrefix(name)
	if err != nil {
		return err
	}
	fmt.Println(core.PrintAccs(ids))
	return nil
}

func AltAcc(cmd cli.Cmd) error {
	var desc *string = nil
	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "desc":
			desc = &fl.V
		}
	}

	if _, err := core.AltAcc(cmd.Pa[0], desc); err != nil {
		return err
	}
	return nil
}

func AltAccRename(cmd cli.Cmd) error {
	old := cmd.Pa[0]
	neo, err := core.SureName(cmd.Pa[1])
	if err != nil {
		return err
	}

	if _, err := core.AltRenameAcc(old, neo); err != nil {
		return err
	}
	return nil
}


func DelAcc(cmd cli.Cmd) error {
	if _, err := core.DelAcc(cmd.Pa[0]); err != nil {
		return err
	}
	return nil
}
