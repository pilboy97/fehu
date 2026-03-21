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

var ErrWrongTxnPattern = errors.New("wrong txn pattern")
var ErrWrongPeriodPattern = errors.New("wrong duration pattern")
var ErrTxnBalance = errors.New("increase and decrease are not equal")

type Pattern struct {
	Name   string
	Amount *money.Money
}

func ParseTxnPattern(pat string) []Pattern {
	var ret = make([]Pattern, 0)

	rexp := `^(~?((\p{L}|_)(\p{L}|\d|_)*:)*((\p{L}|_)(\p{L}|\d|_)*)(<|>)(-|\+)?\d+(\.\d+)?;)*~?((\p{L}|_)(\p{L}|\d|_)*:)*((\p{L}|_)(\p{L}|\d|_)*)(<|>)(-|\+)?\d+(\.\d+)?$`
	ok, err := regexp.MatchString(rexp, pat)
	if err != nil {
		panic(err)
	}

	if !ok {
		panic(ErrWrongTxnPattern)
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
				panic(ErrWrongTxnPattern)
			}

			name, num = tok[0], tok[1]
		} else {
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
			panic(err)
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

	if balance != 0 {
		println(balance)
		panic(ErrTxnBalance)
	}

	return ret
}

func NewTxn(cmd cli.Cmd) {
	var desc string
	var timestamp int64
	var pat = ParseTxnPattern(cmd.Pa[0])

	timestamp = time.Now().UTC().Unix()

	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "desc":
			desc = fl.V
		case "time":
			timestamp = core.ParseTime(fl.V)
		}
	}

	txnID := core.NewTxn(desc, timestamp)
	for _, p := range pat {
		aid := core.GetAccByName(p.Name)
		if aid == -1 {
			core.DelTxn(txnID) // 롤백
			panic(fmt.Errorf("account '%s' not found", p.Name))
		}
		core.NewRecord(txnID, aid, p.Amount)
	}

	fmt.Printf("txn #%d created\n", txnID)
}
func GetTxn(cmd cli.Cmd) {
	var Name string
	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "save":
			Name = core.SureName(fl.V)
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
	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "save":
			Name = core.SureName(fl.V)
		}
	}

	id, err := strconv.ParseInt(cmd.Pa[0], 10, 64)
	if err != nil {
		panic(err)
	}

	fmt.Println(core.PrintTxns([]int64{core.SureID(id)}))

	if len(Name) != 0 {
		core.Vars[Name] = core.NewTable([]int64{id})
	}
}
func GetTxnByDesc(cmd cli.Cmd) {
	var Name string
	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "save":
			Name = core.SureName(fl.V)
		}
	}

	desc := cmd.Pa[0]

	ids := core.GetTxnByDesc(desc)
	fmt.Println(core.PrintTxns(ids))

	if len(Name) != 0 {
		core.Vars[Name] = core.NewTable(ids)
	}
}
func GetTxnByTime(cmd cli.Cmd) {
	var Name string
	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "save":
			Name = core.SureName(fl.V)
		}
	}

	timepat := cmd.Pa[0]
	tokens := strings.Split(timepat, "~")

	if len(tokens) != 2 {
		panic(ErrWrongPeriodPattern)
	}

	var A, B *int64

	if len(tokens[0]) != 0 {
		ts := core.ParseTime(tokens[0])
		A = &ts
	}
	if len(tokens[1]) != 0 {
		ts := core.ParseTime(tokens[1])
		B = &ts
	}

	var ret = core.GetTxnByTime(A, B)
	fmt.Println(core.PrintTxns(ret))
	if len(Name) != 0 {
		core.Vars[Name] = core.NewTable(ret)
	}
}
func AltTxn(cmd cli.Cmd) {
	id, err := strconv.ParseInt(cmd.Pa[0], 10, 64)
	if err != nil {
		panic(err)
	}
	core.SureID(id)
	var desc *string = nil
	var timestamp *int64 = nil
	for i, fl := range cmd.Fl {
		switch fl.F.Name {
		case "desc":
			desc = &cmd.Fl[i].V
		case "time":
			ts := core.ParseTime(fl.V)
			timestamp = &ts
		}
	}

	core.SureID(core.AltTxn(id, desc, timestamp))
}
func AltTxnRecord(cmd cli.Cmd) {
	id, err := strconv.ParseInt(cmd.Pa[0], 10, 64)
	if err != nil {
		panic(err)
	}
	core.SureID(id)

	pat := cmd.Pa[1]
	pats := ParseTxnPattern(pat)

	records := make([]core.Record, len(pats))

	for i, p := range pats {
		aid := core.GetAccByName(p.Name)
		amount := p.Amount

		records[i] = core.Record{
			TID:    id,
			AID:    aid,
			Amount: amount,
		}
	}

	core.SureID(core.AltTxnRecord(id, records))
}
func DelTxn(cmd cli.Cmd) {
	id, err := strconv.ParseInt(cmd.Pa[0], 10, 64)
	if err != nil {
		panic(err)
	}

	core.DelTxn(core.SureID(id))
}
