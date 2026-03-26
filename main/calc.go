package main

import (
	"cli"
	"core"
	"fmt"
	"strings"
)

func Calc(cmd cli.Cmd) {
	stmt := cmd.Pa[0]

	fmt.Println(core.CalcStmt(stmt).String())
}
func Def(cmd cli.Cmd) {
	name := cmd.Pa[0]
	stmt := strings.Join(cmd.Pa[1:], "")

	sureName, err := core.SureName(name)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if err := core.DefStmt(sureName, stmt); err != nil {
		fmt.Println("Error:", err)
	}
}
