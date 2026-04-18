package main

import (
	"cli"
	"core"
	"fmt"
	"strings"
)

func Calc(cmd cli.Cmd) error {
	stmt := strings.Join(cmd.Pa, " ")
	fmt.Println(core.CalcStmt(stmt).String())
	return nil
}

func Def(cmd cli.Cmd) error {
	name := cmd.Pa[0]
	stmt := strings.Join(cmd.Pa[1:], "")

	sureName, err := core.SureName(name)
	if err != nil {
		return err
	}
	if err := core.DefStmt(sureName, stmt); err != nil {
		return err
	}
	return nil
}
