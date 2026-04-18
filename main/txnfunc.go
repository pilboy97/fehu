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
			return nil, ErrWrongTxnPattern
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
				return nil, ErrWrongTxnPattern
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

func NewTxn(cmd cli.Cmd) error {
	var desc string
	var timestamp int64
	var patterns, err = ParseTxnPattern(cmd.Pa[0])

	if err != nil {
		return fmt.Errorf("error parsing transaction pattern: %w", err)
	}

	timestamp = time.Now().UTC().Unix()

	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "desc":
			desc = fl.V
		case "time":
			parsedTime, err := core.ParseTime(fl.V)
			if err != nil {
				return fmt.Errorf("error parsing time: %w", err)
			}
			timestamp = parsedTime
		}
	}

	txnID, err := core.NewTxn(desc, timestamp)
	if err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}

	for _, p := range patterns {
		aid, err := core.GetAccByName(p.Name)
		if err != nil {
			core.DelTxn(txnID) // Rollback
			return fmt.Errorf("account '%s' not found, transaction rolled back", p.Name)
		}
		if _, err := core.NewRecord(txnID, aid, p.Amount); err != nil {
			core.DelTxn(txnID) // Rollback
			return fmt.Errorf("failed to create record for account '%s', transaction rolled back", p.Name)
		}
	}
	fmt.Printf("Transaction #%d created\n", txnID)
	return nil
}

func GetTxn(cmd cli.Cmd) error {
	var Name string
	var err error

	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "save":
			Name, err = core.SureName(fl.V)
			if err != nil {
				return errors.New("invalid variable name")
			}
		}
	}

	ret := core.GetTxn()
	fmt.Println(core.PrintTxns(ret))

	if len(Name) != 0 {
		core.Vars[Name] = core.NewTable(ret)
	}
	return nil
}

func GetTxnByID(cmd cli.Cmd) error {
	var Name string
	var err error

	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "save":
			Name, err = core.SureName(fl.V)
			if err != nil {
				return errors.New("invalid variable name")
			}
		}
	}

	id, err := strconv.ParseInt(cmd.Pa[0], 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing transaction ID: %w", err)
	}

	fmt.Println(core.PrintTxns([]int64{id}))

	if len(Name) != 0 {
		core.Vars[Name] = core.NewTable([]int64{id})
	}
	return nil
}

func GetTxnByDesc(cmd cli.Cmd) error {
	var Name string
	var err error

	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "save":
			Name, err = core.SureName(fl.V)
			if err != nil {
				return errors.New("invalid variable name")
			}
		}
	}

	desc := cmd.Pa[0]

	ids, err := core.GetTxnByDesc(desc)
	if err != nil {
		return err
	}

	fmt.Println(core.PrintTxns(ids))

	if len(Name) != 0 {
		core.Vars[Name] = core.NewTable(ids)
	}
	return nil
}

func GetTxnByTime(cmd cli.Cmd) error {
	var Name string
	var err error

	for _, fl := range cmd.Fl {
		switch fl.F.Name {
		case "save":
			Name, err = core.SureName(fl.V)
			if err != nil {
				return errors.New("invalid variable name")
			}
		}
	}

	timepat := cmd.Pa[0]
	tokens := strings.Split(timepat, "~")

	if len(tokens) != 2 {
		return ErrWrongPeriodPattern
	}

	var A, B *int64

	if len(tokens[0]) != 0 {
		ts, err := core.ParseTime(tokens[0])
		if err != nil {
			return fmt.Errorf("error parsing start time: %w", err)
		}
		A = &ts
	}
	if len(tokens[1]) != 0 {
		ts, err := core.ParseTime(tokens[1])
		if err != nil {
			return fmt.Errorf("error parsing end time: %w", err)
		}
		B = &ts
	}

	var ret = core.GetTxnByTime(A, B)
	fmt.Println(core.PrintTxns(ret))
	if len(Name) != 0 {
		core.Vars[Name] = core.NewTable(ret)
	}
	return nil
}

func AltTxn(cmd cli.Cmd) error {
	id, err := strconv.ParseInt(cmd.Pa[0], 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing transaction ID: %w", err)
	}

	var desc *string = nil
	var timestamp *int64 = nil
	for i, fl := range cmd.Fl {
		switch fl.F.Name {
		case "desc":
			desc = &cmd.Fl[i].V
		case "time":
			parsedTime, err := core.ParseTime(fl.V)
			if err != nil {
				return fmt.Errorf("error parsing time: %w", err)
			}
			timestamp = &parsedTime
		}
	}

	if _, err := core.AltTxn(id, desc, timestamp); err != nil {
		return fmt.Errorf("error altering transaction: %w", err)
	}
	fmt.Printf("Transaction #%d updated successfully.\n", id)
	return nil
}

func AltTxnRecord(cmd cli.Cmd) error {
	id, err := strconv.ParseInt(cmd.Pa[0], 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing transaction ID: %w", err)
	}

	pat := cmd.Pa[1]
	pats, err := ParseTxnPattern(pat)
	if err != nil {
		return fmt.Errorf("error parsing transaction pattern: %w", err)
	}

	records := make([]core.Record, len(pats))

	for i, p := range pats {
		aid, err := core.GetAccByName(p.Name)
		if err != nil {
			return fmt.Errorf("account '%s' not found", p.Name)
		}

		records[i] = core.Record{
			TID:    id,
			AID:    aid,
			Amount: p.Amount,
		}
	}

	if _, err := core.AltTxnRecord(id, records); err != nil {
		return fmt.Errorf("error altering transaction record: %w", err)
	}
	fmt.Printf("Transaction #%d records updated successfully.\n", id)
	return nil
}

func DelTxn(cmd cli.Cmd) error {
	id, err := strconv.ParseInt(cmd.Pa[0], 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing transaction ID: %w", err)
	}

	if _, err := core.DelTxn(id); err != nil {
		return fmt.Errorf("error deleting transaction: %w", err)
	}
	fmt.Printf("Transaction #%d deleted successfully.\n", id)
	return nil
}
