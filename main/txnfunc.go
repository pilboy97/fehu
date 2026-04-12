package main

import (
	"cli"
	"core"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/pkg/errors"
)

var ErrWrongTxnPattern = errors.New("wrong transaction pattern")
var ErrWrongPeriodPattern = errors.New("wrong duration pattern")
var ErrTxnBalance = errors.New("increase and decrease are not equal")

type Pattern struct {
	Name   string
	Amount *money.Money
}

func ParseTxnPattern(pat string) ([]Pattern, error) {
	var ret = make([]Pattern, 0)

	rexp := `^(~?((\p{L}|_)(\p{L}|\d|_)*:)*((\p{L}|_)(\p{L}|\d|_)*)(<|>)(-|\+)?\d+(\.\d+)?;)*~?((\p{L}|_)(\p{L}|\d|_)*:)*((\p{L}|_)(\p{L}|\d|_)*)(<|>)(-|\+)?\d+(\.\d+)?$`
	ok, err := regexp.MatchString(rexp, pat)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, ErrWrongTxnPattern
	}

	rexp = `^(?P<name>~?((\p{L}|_)(\p{L}|\d|_)*:)*((\p{L}|_)(\p{L}|\d|_)*))(<|>)(?P<num>(-|\+)?\d+(\.\d+)?)$`
	re := regexp.MustCompile(rexp)

	records := strings.Split(pat, ";")

	var balance float64

	for _, ch := range records {

		if !re.MatchString(ch) {
			panic(ErrWrongTxnPattern)
		}

		var dir bool = false
		var name, num string
		var isNeg bool

		if strings.Contains(ch, ">") {
			isNeg = true

			tok := strings.Split(ch, ">")
			if len(tok) != 2 {
				return nil, ErrWrongTxnPattern
			}

			name, num = tok[0], tok[1]
		} else { // Contains "<"
			isNeg = false
			dir = true
			tok := strings.Split(ch, "<")
			if len(tok) != 2 {
				panic(ErrWrongTxnPattern)
			}

			name, num = tok[0], tok[1]
		}

		if strings.HasPrefix(name, "~") {
			isNeg = !isNeg
		}

		f, err := strconv.ParseFloat(num, 64)
		if err != nil {
			return nil, err
		}

		if dir {
			balance += f
		} else {
			balance -= f
		}

		if isNeg {
			f = -f
		}

		amount := money.NewFromFloat(f, env.Code())

		ret = append(ret, Pattern{Name: name, Amount: amount})
	}

	if balance != 0 { // Check if the transaction is balanced
		return nil, ErrTxnBalance
	}
	return ret, nil
}

func NewTxn(cmd cli.Cmd) {
	var desc string
	var timestamp int64
	var patterns, err = ParseTxnPattern(cmd.Pa[0])

	if err != nil {
		fmt.Println("Error parsing transaction pattern:", err)
		return
	}

	timestamp = time.Now().UTC().Unix()

	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "desc":
			desc = fl.V
		case "time": // TODO: Handle error from ParseTime
			parsedTime, err := core.ParseTime(fl.V)
			if err != nil {
				fmt.Println("Error parsing time:", err)
				return
			}
			timestamp = parsedTime
		}
	}

	txnID, err := core.NewTxn(desc, timestamp)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}

	for _, p := range patterns {
		aid, err := core.GetAccByName(p.Name)

		if err != nil {
			core.DelTxn(txnID) // Rollback
			fmt.Printf("Error: account '%s' not found. Transaction rolled back.\n", p.Name)
			return
		}
		if _, err := core.NewRecord(txnID, aid, p.Amount); err != nil {
			core.DelTxn(txnID) // Rollback
			fmt.Printf("Error: failed to create record for account '%s'. Transaction rolled back.\n", p.Name)
			return
		}
	}
	fmt.Printf("Transaction #%d created\n", txnID)
}
func GetTxn(cmd cli.Cmd) {
	var Name string
	var err error

	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "save":
			Name, err = core.SureName(fl.V)
			if err != nil {
				fmt.Println("Invalid variable name")
				return
			}
		}
	}

	ret := core.GetTxn()
	fmt.Println(core.PrintTxns(ret))

	if len(Name) != 0 {
		core.Vars[Name] = core.NewTable(ret)
	}
}
func GetTxnByID(cmd cli.Cmd) {
	var Name string
	var err error

	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "save":
			Name, err = core.SureName(fl.V)
			if err != nil {
				fmt.Println("Invalid variable name")
				return
			}
		}
	}

	id, err := strconv.ParseInt(cmd.Pa[0], 10, 64)
	if err != nil {
		fmt.Println("Error parsing transaction ID:", err)
	}

	fmt.Println(core.PrintTxns([]int64{id}))

	if len(Name) != 0 {
		core.Vars[Name] = core.NewTable([]int64{id})
	}
}
func GetTxnByDesc(cmd cli.Cmd) {
	var Name string
	var err error

	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "save":
			Name, err = core.SureName(fl.V)
			if err != nil {
				fmt.Println("Invalid variable name")
				return
			}
		}
	}

	desc := cmd.Pa[0]

	ids, err := core.GetTxnByDesc(desc)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(core.PrintTxns(ids))

	if len(Name) != 0 {
		core.Vars[Name] = core.NewTable(ids)
	}
}
func GetTxnByTime(cmd cli.Cmd) {
	var Name string
	var err error

	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "save":
			Name, err = core.SureName(fl.V)
			if err != nil {
				fmt.Println("Invalid variable name")
				return
			}
		}
	}

	timepat := cmd.Pa[0]
	tokens := strings.Split(timepat, "~")

	if len(tokens) != 2 {
		fmt.Println("Error:", ErrWrongPeriodPattern)
	}

	var A, B *int64

	if len(tokens[0]) != 0 {
		ts, err := core.ParseTime(tokens[0])
		if err != nil {
			fmt.Println("Error parsing start time:", err)
			return
		}

		A = &ts
		if err != nil {
			fmt.Println("Error parsing start time:", err)
			return
		}
	}
	if len(tokens[1]) != 0 {
		ts, err := core.ParseTime(tokens[1])
		if err != nil {
			fmt.Println("Error parsing start time:", err)
			return
		}

		B = &ts
		if err != nil {
			fmt.Println("Error parsing end time:", err)
			return
		}
	} // TODO: Handle error from ParseTime

	var ret = core.GetTxnByTime(A, B)
	fmt.Println(core.PrintTxns(ret))
	if len(Name) != 0 {
		core.Vars[Name] = core.NewTable(ret)
	}
}
func AltTxn(cmd cli.Cmd) {
	id, err := strconv.ParseInt(cmd.Pa[0], 10, 64)
	if err != nil {
		fmt.Println("Error parsing transaction ID:", err)
		return
	}

	var desc *string = nil
	var timestamp *int64 = nil
	for i, fl := range cmd.Fl {
		switch fl.F.Name {
		case "desc":
			desc = &cmd.Fl[i].V
		case "time": // TODO: Handle error from ParseTime
			parsedTime, err := core.ParseTime(fl.V)
			if err != nil {
				fmt.Println("Error parsing time:", err)
				return
			}
			timestamp = &parsedTime
		}
	}

	if _, err := core.AltTxn(id, desc, timestamp); err != nil {
		fmt.Println("Error altering transaction:", err)
		return
	}
	fmt.Printf("Transaction #%d updated successfully.\n", id)
}
func AltTxnRecord(cmd cli.Cmd) {
	id, err := strconv.ParseInt(cmd.Pa[0], 10, 64)
	if err != nil {
		fmt.Println("Error parsing transaction ID:", err)
		return
	}

	pat := cmd.Pa[1]
	pats, err := ParseTxnPattern(pat)
	if err != nil {
		fmt.Println("Error parsing transaction pattern:", err)
		return
	}

	records := make([]core.Record, len(pats))

	for i, p := range pats { // TODO: Handle error from ParseTxnPattern
		aid, err := core.GetAccByName(p.Name)
		if err != nil {
			fmt.Printf("Error: account '%s' not found.\n", p.Name)
			return
		}

		amount := p.Amount

		records[i] = core.Record{
			TID:    id,
			AID:    aid,
			Amount: amount,
		} // TODO: Handle error from GetAccByName
	}

	if _, err := core.AltTxnRecord(id, records); err != nil {
		fmt.Println("Error altering transaction record:", err)
		return
	}
	fmt.Printf("Transaction #%d records updated successfully.\n", id)
}
func DelTxn(cmd cli.Cmd) {
	id, err := strconv.ParseInt(cmd.Pa[0], 10, 64)
	if err != nil {
		fmt.Println("Error parsing transaction ID:", err)
		return
	}

	if _, err := core.DelTxn(id); err != nil {
		fmt.Println("Error deleting transaction:", err)
		return
	}
	fmt.Printf("Transaction #%d deleted successfully.\n", id)
}
