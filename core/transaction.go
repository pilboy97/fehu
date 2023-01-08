package core

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func NewTxn(desc string, timestamp time.Time) int64 {
	ChkDB()

	stmt := `insert into txn(desc, time) values(?,?)`
	res, err := DB.Exec(stmt, desc, timestamp)
	if err != nil {
		panic(err)
	}

	ret, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	CreateTagInDesc(desc, nil, []int64{ret})

	return ret
}
func GetTxn() []int64 {
	ChkDB()

	var ret []int64

	stmt := `select id from txn order by time`
	rows, err := DB.Query(stmt)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id int64
		rows.Scan(&id)

		ret = append(ret, id)
	}

	return ret
}
func GetTxnByID(tid int64) Txn {
	ChkDB()

	var ret Txn

	stmt := `select id, desc, time from txn where id=? order by time`
	err := DB.QueryRow(stmt, tid).Scan(
		&ret.ID,
		&ret.Desc,
		&ret.Time,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		ret.ID = -1
		return ret
	}

	records := GetRecordByTID(tid)
	ret.Record = make([]Record, len(records))

	for i, rid := range records {
		r := GetRecordByID(rid)
		ret.Record[i] = r
	}

	return ret
}
func GetTxnByTime(st *time.Time, ed *time.Time) []int64 {
	ChkDB()

	var ret []int64 = make([]int64, 0)
	var rows *sql.Rows
	var err error

	if st != nil && ed != nil && st.After(*ed) {
		st, ed = ed, st
	}

	switch {
	case st != nil && ed != nil:
		rows, err = DB.Query(`select id from txn where time between ? and ? order by time`, st, ed)
		if err != nil {
			panic(err)
		}
	case st != nil:
		rows, err = DB.Query(`select id from txn where time > ? order by time`, st)
		if err != nil {
			panic(err)
		}
	case ed != nil:
		rows, err = DB.Query(`select id from txn where time < ? order by time`, ed)
		if err != nil {
			panic(err)
		}
	default:
		rows, err = DB.Query(`select id from txn order by time`)
		if err != nil {
			panic(err)
		}
	}

	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			panic(err)
		}

		ret = append(ret, id)
	}

	return ret
}
func GetTxnByDesc(desc string) []int64 {
	var ret []int64

	stmt := `select id from txn where instr(desc,?) > 0 order by time`
	row, err := DB.Query(stmt, desc)
	if err != nil {
		panic(err)
	}

	ret = make([]int64, 0)
	for row.Next() {
		var id int64
		err = row.Scan(&id)
		if err != nil {
			panic(err)
		}

		ret = append(ret, id)
	}

	return ret
}
func AltTxn(tid int64, desc *string, timestamp *time.Time) int64 {
	if GetTxnByID(tid).ID == -1 {
		return -1
	}
	if desc != nil {
		stmt := `update Txn set desc=? where id=?`
		_, err := DB.Exec(stmt, desc, tid)
		if err != nil {
			panic(err)
		}
	}
	if timestamp != nil {
		stmt := `update Txn set time=? where id=?`
		_, err := DB.Exec(stmt, timestamp, tid)
		if err != nil {
			panic(err)
		}
	}

	CreateTagInDesc(GetAccByID(tid).Desc, nil, []int64{tid})
	return tid
}
func AltTxnRecord(tid int64, pat []Record) int64 {
	if GetTxnByID(tid).ID == -1 {
		return -1
	}

	for _, p := range pat {
		if p.TID != tid || p.AID == -1 {
			return -1
		}
	}

	old := GetRecordByTID(tid)

	for _, id := range old {
		DelRecord(id)
	}

	for _, p := range pat {
		aid := p.AID
		amount := p.Amount

		NewRecord(tid, aid, amount)
	}

	return tid
}
func DelTxn(tid int64) int64 {
	if GetTxnByID(tid).ID == -1 {
		return -1
	}

	stmt := `delete from Txn where id=?`
	_, err := DB.Exec(stmt, tid)
	if err != nil {
		panic(err)
	}

	return tid
}

func PrintRecord(id int64) string {
	SureID(id)

	var ch = GetRecordByID(id)
	var acc = GetAccByID(ch.AID)
	var amount = ch.Amount

	return fmt.Sprintf("        |                        |%24s|%8s|        ", acc.Name, amount.Display())
}
func PrintTxn(id int64) string {
	SureID(id)

	var ret = make([]string, 0)

	var txn = GetTxnByID(id)
	var records = GetRecordByTID(id)

	ret = append(ret, fmt.Sprintf("%8d|%24s|                        |        |%16s", txn.ID, txn.Time.Format(TimeFmt), txn.Desc))
	for _, ch := range records {
		ret = append(ret, PrintRecord(ch))
	}
	return strings.Join(ret, "\n")
}
func PrintTxns(txn []int64) string {
	var ret = make([]string, 0)
	ret = append(ret, "      id|                    time|                    name|  amount|            desc")
	for _, id := range txn {
		ret = append(ret, PrintTxn(id))
	}
	ret = append(ret, fmt.Sprintf("%8d\ttxn(s)\tfound", len(txn)))

	return strings.Join(ret, "\n")
}
