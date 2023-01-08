package main

import (
	"cli"
	"core"
	"fmt"
)

func Calc(cmd cli.Cmd) {
	stmt := cmd.Pa[0]

	fmt.Println(core.CalcStmt(stmt).String())
}
func Def(cmd cli.Cmd) {
	name := cmd.Pa[0]
	stmt := cmd.Pa[1]

	core.SureName(name)
	core.DefStmt(name, stmt)
}
