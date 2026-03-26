package main

import (
	"cli"
	"core"
	"fmt"
	"log"

	"github.com/Rhymond/go-money"
)

func Proc(cmd cli.Cmd) error {
	if cmd.St.Param > len(cmd.Pa) || len(cmd.St.Flags) < len(cmd.Fl) {
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
	case altCodeState:
		AltCode(cmd)
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

	if err := core.Open(path + ".db"); err != nil {
		panic(err)
	}
}
func Close(cmd cli.Cmd) {
	core.Close()
}
func AltCode(cmd cli.Cmd) {
	code := cmd.Pa[0]
	core.Code = code
	if money.GetCurrency(core.Code) == nil {
		money.AddCurrency(core.Code, core.Code+" ", "1 $", ".", ",", 8)
	}

	fmt.Printf("Currency code changed to %s\n", core.Code)

	if err := core.SetConfig(core.Config{
		Currency: core.Code,
	}); err != nil {
		fmt.Printf("Error saving currency to config: %v\n", err)
		return
	}
}
