package core

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Rhymond/go-money"
)

func NewAcc(name string, desc string) int64 {
	ChkDB()

	stmt := `insert into acc(name, desc) values(?,?)`
	res, err := DB.Exec(stmt, name, desc)
	if err != nil {
		panic(err)
	}

	ret, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	CreateTagInDesc(desc, []int64{ret}, nil)

	return ret
}
func GetAcc() []int64 {
	ChkDB()

	var ret []int64
	stmt := `select id from acc order by id`
	row, err := DB.Query(stmt)
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
func GetAccByID(aid int64) Acc {
	ChkDB()

	var ret Acc
	stmt := `select id, name, desc from acc where id=? order by id`
	err := DB.QueryRow(stmt, aid).Scan(
		&ret.ID,
		&ret.Name,
		&ret.Desc,
	)
	if err != nil {
		panic(err)
	}
	return ret
}
func GetAccByName(name string) int64 {
	ChkDB()

	var ret int64

	stmt := `select id from acc where name=? order by id`
	err := DB.QueryRow(stmt, name).Scan(&ret)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		return -1
	}

	return ret
}
func GetAccByPrefix(name string) []int64 {
	ChkDB()

	var ret []int64

	name = name + "%"

	stmt := `select id from acc where name like ? order by id`
	row, err := DB.Query(stmt, name)
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
func GetAccByDesc(desc string) []int64 {
	ChkDB()

	var ret []int64

	stmt := `select id from acc where instr(desc,?) > 0 order by id`
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
func GetAccAmount(id int64) (*money.Money, error) {
	ChkDB()

	if GetAccByID(id).ID == -1 {
		return &money.Money{}, nil
	}

	var ret *money.Money
	ret = money.New(0, Code)

	records := GetRecordByAID(id)
	for _, rid := range records {
		record := GetRecordByID(rid)

		var err error
		ret, err = ret.Add(record.Amount)
		if err != nil {
			return &money.Money{}, err
		}
	}

	return ret, nil
}
func AltAcc(name string, desc *string) int64 {
	ChkDB()

	ID := GetAccByName(name)
	if ID == -1 {
		return -1
	}

	if desc == nil {
		return ID
	}

	if desc != nil {
		stmt := `update acc set desc=? where id=?`
		_, err := DB.Exec(stmt, desc, ID)
		if err != nil {
			panic(err)
		}
	}

	CreateTagInDesc(GetAccByID(ID).Desc, []int64{ID}, nil)

	return ID
}
func AltRenameAcc(old, neo string) int64 {
	ChkDB()

	ID := GetAccByName(old)
	if ID == -1 {
		return -1
	}

	if ID2 := GetAccByName(neo); ID2 != -1 {
		return -2
	}

	stmt := `update acc set name=? where id=?`
	_, err := DB.Exec(stmt, neo, ID)
	if err != nil {
		panic(err)
	}

	return ID
}
func DelAcc(name string) int64 {
	ChkDB()

	ID := GetAccByName(name)
	if ID == -1 {
		return -1
	}

	stmt := `delete from acc where id=?`
	_, err := DB.Exec(stmt, ID)
	if err != nil {
		panic(err)
	}

	return ID
}

func PrintAccs(acc []int64) string {
	var ret = []string{}

	ret = append(ret, "      id|                    name|  amount|                    desc")
	for _, id := range acc {
		var acc = GetAccByID(id)
		amount, err := GetAccAmount(id)
		if err != nil {
			panic(err)
		}

		ret = append(ret, fmt.Sprintf("%8d|%24s|%8s|%24s", acc.ID, acc.Name, amount.Display(), acc.Desc))
	}
	ret = append(ret, fmt.Sprintf("%8d\tacc(s)\tfound", len(acc)))

	return strings.Join(ret, "\n")
}
