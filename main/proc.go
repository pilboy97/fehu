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
		return Open(cmd)
	case closeState:
		Close(cmd)
	case quitState:
		return cli.ErrShutdownSystem
	case calcState:
		return Calc(cmd)
	case defState:
		return Def(cmd)

	case newState:
		fmt.Println(cmd.St.HelpString())
	case getState:
		fmt.Println(cmd.St.HelpString())
	case altState:
		fmt.Println(cmd.St.HelpString())
	case delState:
		fmt.Println(cmd.St.HelpString())

	case newAccState:
		return NewAcc(cmd)
	case getAccState:
		return GetAcc(cmd)
	case getAccByNameState:
		return GetAccByName(cmd)
	case getAccByDescState:
		return GetAccByDesc(cmd)
	case getAccChildState:
		return GetAccChild(cmd)
	case altAccState:
		return AltAcc(cmd)
	case altAccRenameState:
		return AltAccRename(cmd)
	case altCodeState:
		return AltCode(cmd)
	case delAccState:
		return DelAcc(cmd)

	case newTxnState:
		return NewTxn(cmd)

	case getTxnState:
		return GetTxn(cmd)
	case getTxnByIDState:
		return GetTxnByID(cmd)
	case getTxnByDescState:
		return GetTxnByDesc(cmd)
	case getTxnByTimeState:
		return GetTxnByTime(cmd)

	case altTxnState:
		return AltTxn(cmd)
	case altTxnRecordState:
		return AltTxnRecord(cmd)
	case delTxnState:
		return DelTxn(cmd)

	case newTagState:
		return NewTag(cmd)
	case getTagState:
		return GetTag(cmd)
	case getTagByNameState:
		return GetTagByName(cmd)
	case getTagByDescState:
		return GetTagByDesc(cmd)
	case altTagState:
		return AltTag(cmd)
	case altTagRenameState:
		return AltTagRename(cmd)
	case delTagState:
		return DelTag(cmd)

	default:
		log.Printf("%s is not supported yet", cmd.St.Name)
	}

	return nil
}

func Open(cmd cli.Cmd) error {
	path := cmd.Pa[0]
	if err := core.Open(path + ".db"); err != nil {
		return err
	}
	return nil
}

func Close(cmd cli.Cmd) {
	core.Close()
}

func AltCode(cmd cli.Cmd) error {
	code := cmd.Pa[0]
	core.Code = code
	if money.GetCurrency(core.Code) == nil {
		money.AddCurrency(core.Code, core.Code+" ", "1 $", ".", ",", 8)
	}

	fmt.Printf("Currency code changed to %s\n", core.Code)

	if err := core.SetConfig(core.Config{
		Currency: core.Code,
	}); err != nil {
		return fmt.Errorf("error saving currency to config: %w", err)
	}
	return nil
}
