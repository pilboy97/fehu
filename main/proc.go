package main

import (
	"cli"
	"core"
	"fmt"
	"log"
)

func Proc(cmd cli.Cmd) error {
	if cmd.St.Param != len(cmd.Pa) || len(cmd.St.Flags) < len(cmd.Fl) {
		fmt.Println(cmd.St.HelpString())
		return nil
	}

	switch cmd.St {
	case rootState:
		fmt.Println(cmd.St.HelpString())

	case openState:
		Open(cmd)
	case closeState:
		Close(cmd)
	case quitState:
		return cli.ErrShutdownSystem
	case calcState:
		Calc(cmd)
	case defState:
		Def(cmd)

	case newState:
		fmt.Println(cmd.St.HelpString())
	case getState:
		fmt.Println(cmd.St.HelpString())
	case altState:
		fmt.Println(cmd.St.HelpString())
	case delState:
		fmt.Println(cmd.St.HelpString())

	case newAccState:
		NewAcc(cmd)
	case getAccState:
		GetAcc(cmd)
	case getAccByNameState:
		GetAccByName(cmd)
	case getAccByDescState:
		GetAccByDesc(cmd)
	case getAccChildState:
		GetAccChild(cmd)
	case altAccState:
		AltAcc(cmd)
	case altAccRenameState:
		AltAccRename(cmd)
	case delAccState:
		DelAcc(cmd)

	case newTxnState:
		NewTxn(cmd)

	case getTxnState:
		GetTxn(cmd)
	case getTxnByIDState:
		GetTxnByID(cmd)
	case getTxnByDescState:
		GetTxnByDesc(cmd)
	case getTxnByTimeState:
		GetTxnByTime(cmd)

	case altTxnState:
		AltTxn(cmd)
	case altTxnRecordState:
		AltTxnRecord(cmd)
	case delTxnState:
		DelTxn(cmd)

	case newTagState:
		NewTag(cmd)
	case getTagState:
		GetTag(cmd)
	case getTagByNameState:
		GetTagByName(cmd)
	case getTagByDescState:
		GetTagByDesc(cmd)
	case altTagState:
		AltTag(cmd)
	case altTagRenameState:
		AltTagRename(cmd)
	case delTagState:
		DelTag(cmd)

	default:
		log.Printf("%s is not supported yet", cmd.St.Name)
	}

	return nil
}
func Open(cmd cli.Cmd) {
	path := cmd.Pa[0]
	core.Open(path + ".db")
}
func Close(cmd cli.Cmd) {
	core.Close()
}
