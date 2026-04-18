package main

import (
	"cli"
	"core"
	"fmt"
)

func NewTag(cmd cli.Cmd) error {
	var tag core.Tag
	var err error
	tag.Name, err = core.SureName(cmd.Pa[0])
	if err != nil {
		return err
	}

	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "desc":
			tag.Desc = fl.V
		}
	}

	id, err := core.NewTag(tag.Name, tag.Desc)
	if err != nil {
		return err
	}
	fmt.Printf("tag %d created\n", id)
	return nil
}

func GetTag(cmd cli.Cmd) error {
	ret, err := core.GetTag()
	if err != nil {
		return err
	}
	fmt.Println(core.PrintTags(ret))
	return nil
}

func GetTagByName(cmd cli.Cmd) error {
	id, err := core.GetTagByName(cmd.Pa[0])
	if err != nil {
		return err
	}
	fmt.Println(core.PrintTags([]int64{id}))
	return nil
}

func GetTagByDesc(cmd cli.Cmd) error {
	ids, err := core.GetTagByDesc(cmd.Pa[0])
	if err != nil {
		return err
	}
	fmt.Println(core.PrintTags(ids))
	return nil
}

func AltTag(cmd cli.Cmd) error {
	var desc *string = nil
	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "desc":
			desc = &fl.V
		}
	}

	if _, err := core.AltTag(cmd.Pa[0], desc); err != nil {
		return err
	}
	return nil
}

func AltTagRename(cmd cli.Cmd) error {
	old := cmd.Pa[0]
	neo, err := core.SureName(cmd.Pa[1])
	if err != nil {
		return err
	}

	if _, err := core.AltRenameTag(old, neo); err != nil {
		return err
	}
	return nil
}

func DelTag(cmd cli.Cmd) error {
	if _, err := core.DelTag(cmd.Pa[0]); err != nil {
		return err
	}
	return nil
}
