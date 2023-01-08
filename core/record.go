package core

import (
	"database/sql"

	"github.com/Rhymond/go-money"
	"github.com/pkg/errors"
)

var ErrUnexpectedTID = errors.New("unexpected tid appears")

func NewRecord(tid int64, aid int64, amount *money.Money) int64 {
	ChkDB()

	stmt := `insert into record(tid, aid, amount) values(?,?,?)`
	res, err := DB.Exec(stmt, tid, aid, amount.Amount())
	if err != nil {
		panic(err)
	}

	ret, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	return ret
}
func GetRecordByID(id int64) Record {
	ChkDB()

	var ret Record
	var raw int64
	stmt := `select id, tid, aid, amount from record where id=?`
	err := DB.QueryRow(stmt, id).Scan(
		&ret.ID,
		&ret.TID,
		&ret.AID,
		&raw,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		ret.ID = -1
		return ret
	}

	ret.Amount = money.New(raw, Code)
	return ret
}
func GetRecordByTAID(tid, aid int64) Record {
	ChkDB()

	var ret Record
	var raw int64
	stmt := `select id, tid, aid, amount from record where tid=? and aid=?`
	err := DB.QueryRow(stmt, tid, aid).Scan(
		&ret.ID,
		&ret.TID,
		&ret.AID,
		&raw,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		ret.ID = -1
		return ret
	}

	ret.Amount = money.New(raw, Code)

	return ret
}
func GetRecordByTID(tid int64) []int64 {
	ChkDB()

	var ret []int64

	stmt := `select id from record where tid=?`
	row, err := DB.Query(stmt, tid)
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
func GetRecordByAID(aid int64) []int64 {
	ChkDB()

	var ret []int64

	stmt := `select id from record where aid=?`
	row, err := DB.Query(stmt, aid)
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
func AltRecord(id int64, tid *int64, aid *int64, amount *money.Money) int64 {
	ChkDB()

	if tid != nil {
		stmt := `update record set tid=? where id=?`
		_, err := DB.Exec(stmt, tid, id)
		if err != nil {
			panic(err)
		}
	}
	if aid != nil {
		stmt := `update record set aid=? where id=?`
		_, err := DB.Exec(stmt, aid, id)
		if err != nil {
			panic(err)
		}
	}
	if amount != nil {
		raw := amount.Amount()
		stmt := `update record set amount=? where id=?`
		_, err := DB.Exec(stmt, raw, id)
		if err != nil {
			panic(err)
		}
	}

	return id
}
func DelRecord(id int64) int64 {
	ChkDB()

	stmt := `delete from record where id=?`
	_, err := DB.Exec(stmt, id)
	if err != nil {
		panic(err)
	}

	return id
}
