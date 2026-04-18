package core

import (
	"database/sql"
	"fmt"

	"github.com/Rhymond/go-money"
	"github.com/pkg/errors"
)

var ErrUnexpectedTID = errors.New("unexpected tid appears")

func NewRecord(tid int64, aid int64, amount *money.Money) (int64, error) {
	MustDB()

	stmt := `insert into record(tid, aid, amount) values(?,?,?)`
	res, err := DB.Exec(stmt, tid, aid, amount.Amount())

	if err != nil {
		return 0, err
	}

	ret, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return ret, nil
}
func GetRecordByID(id int64) (Record, error) {
	MustDB()

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
			return Record{}, err
		}
		return Record{ID: -1}, ErrCannotFind(fmt.Sprintf("record with ID %d", id))
	}

	ret.Amount = money.New(raw, Code)
	return ret, nil
}
func GetRecordByTAID(tid, aid int64) Record {
	MustDB()

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
func GetRecordByTID(tid int64) ([]int64, error) {
	MustDB()

	rows, err := DB.Query(`select id from record where tid=?`, tid)
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

func GetRecordByAID(aid int64) ([]int64, error) {
	MustDB()

	rows, err := DB.Query(`select id from record where aid=?`, aid)
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
func AltRecord(id int64, tid *int64, aid *int64, amount *money.Money) (int64, error) {
	MustDB()

	if tid != nil {
		stmt := `update record set tid=? where id=?`
		_, err := DB.Exec(stmt, tid, id)
		if err != nil {
			return 0, err
		}
	}
	if aid != nil {
		stmt := `update record set aid=? where id=?`
		_, err := DB.Exec(stmt, aid, id)
		if err != nil {
			return 0, err
		}
	}
	if amount != nil {
		raw := amount.Amount()
		stmt := `update record set amount=? where id=?`
		_, err := DB.Exec(stmt, raw, id)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}
func DelRecord(id int64) (int64, error) {
	MustDB()

	stmt := `delete from record where id=?`
	_, err := DB.Exec(stmt, id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
