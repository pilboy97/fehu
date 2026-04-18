package core

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func NewTxn(desc string, timestamp int64) (int64, error) {
	MustDB()

	stmt := `insert into txn(desc, time) values(?,?)`
	res, err := DB.Exec(stmt, desc, timestamp)
	if err != nil {
		return 0, err
	}

	ret, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	CreateTagInDesc(desc, nil, []int64{ret})

	return ret, nil
}
func GetTxn() []int64 {
	MustDB()

	var ret []int64

	stmt := `select id from txn order by time`
	rows, err := DB.Query(stmt)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			panic(err)
		}
		ret = append(ret, id)
	}

	return ret
}
func GetTxnByID(tid int64) (Txn, error) {
	MustDB()

	var ret Txn

	stmt := `select id, desc, time from txn where id=?`
	err := DB.QueryRow(stmt, tid).Scan(
		&ret.ID,
		&ret.Desc,
		&ret.Time,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			return Txn{}, err
		}
		return Txn{ID: -1}, ErrCannotFind(fmt.Sprintf("transaction with ID %d", tid))
	}

	records, err := GetRecordByTID(tid)
	if err != nil {
		return Txn{}, err
	}
	ret.Record = make([]Record, len(records))

	for i, rid := range records {
		r, err := GetRecordByID(rid)
		if err != nil {
			return Txn{}, err
		}
		ret.Record[i] = r
	}

	return ret, nil
}
func GetTxnByTime(st *int64, ed *int64) []int64 {
	MustDB()

	var ret []int64 = make([]int64, 0)
	var rows *sql.Rows
	var err error

	if st != nil && ed != nil && *st > *ed {
		*st, *ed = *ed, *st
	}

	switch {
	case st != nil && ed != nil:
		rows, err = DB.Query(`select id from txn where time between ? and ? order by time`, st, ed)
	case st != nil:
		rows, err = DB.Query(`select id from txn where time > ? order by time`, st)
	case ed != nil:
		rows, err = DB.Query(`select id from txn where time < ? order by time`, ed)
	default:
		rows, err = DB.Query(`select id from txn order by time`)
	}
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			panic(err)
		}
		ret = append(ret, id)
	}

	return ret
}
func GetTxnByDesc(desc string) ([]int64, error) {
	MustDB()

	stmt := `select id from txn where instr(desc,?) > 0 order by time`
	rows, err := DB.Query(stmt, desc)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ret := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		ret = append(ret, id)
	}

	return ret, nil
}
func AltTxn(tid int64, desc *string, timestamp *int64) (int64, error) {
	MustDB()

	txn, err := GetTxnByID(tid)
	if err != nil {
		return 0, err
	}
	if txn.ID == -1 {
		return 0, ErrCannotFind(fmt.Sprintf("transaction with ID %d", tid))
	}
	if desc != nil {
		_, err := DB.Exec(`delete from tagtxn where tid=?`, tid)
		if err != nil {
			return 0, err
		}

		stmt := `update txn set desc=? where id=?`
		_, err = DB.Exec(stmt, desc, tid)
		if err != nil {
			return 0, err
		}
	}
	if timestamp != nil {
		stmt := `update txn set time=? where id=?`
		_, err := DB.Exec(stmt, timestamp, tid)
		if err != nil {
			return 0, err
		}
	}

	if desc != nil {
		CreateTagInDesc(*desc, nil, []int64{tid})

		if err := CleanUpTags(); err != nil {
			return 0, err
		}
	}
	return tid, nil
}
func AltTxnRecord(tid int64, pat []Record) (int64, error) {
	MustDB()

	txn, err := GetTxnByID(tid)
	if err != nil {
		return 0, err
	}
	if txn.ID == -1 {
		return 0, ErrCannotFind(fmt.Sprintf("transaction with ID %d", tid))
	}

	for _, p := range pat {
		if p.TID != tid || p.AID == -1 {
			return 0, fmt.Errorf("invalid record pattern for transaction %d", tid)
		}
	}

	old, err := GetRecordByTID(tid)
	if err != nil {
		return 0, err
	}

	for _, id := range old {
		if _, err := DelRecord(id); err != nil {
			return 0, err
		}
	}

	for _, p := range pat {
		aid := p.AID
		amount := p.Amount

		if _, err := NewRecord(tid, aid, amount); err != nil {
			return 0, err
		}
	}

	return tid, nil
}
func DelTxn(tid int64) (int64, error) {
	MustDB()

	txn, err := GetTxnByID(tid)
	if err != nil {
		return 0, err
	}
	if txn.ID == -1 {
		return 0, ErrCannotFind(fmt.Sprintf("transaction with ID %d", tid))
	}

	stmt := `delete from txn where id=?`
	_, err = DB.Exec(stmt, tid)
	if err != nil {
		return 0, err
	}

	if err := CleanUpTags(); err != nil {
		return 0, err
	}
	return tid, nil
}

func PrintRecord(id int64) string {
	ch, err := GetRecordByID(id)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	acc, err := GetAccByID(ch.AID)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	var amount = ch.Amount.Absolute()

	return fmt.Sprintf("        |                        |%24s|%8s|        ", acc.Name, amount.Absolute().Display())
}
func PrintTxn(id int64) string {
	var ret = make([]string, 0)

	var txn, err = GetTxnByID(id)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	records, err := GetRecordByTID(id)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	// UTC 타임스탬프를 로컬 시간으로 변환하여 출력
	t := time.Unix(txn.Time, 0).Local()

	ret = append(ret, fmt.Sprintf("%8d|%24s|                        |        |%16s", txn.ID, t.Format(TimeFmt), txn.Desc))
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
