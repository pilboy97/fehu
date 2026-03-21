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
	stmt := `select id, name, desc from acc where id=?`
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

	stmt := `select id from acc where name=?`
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
		_, err := DB.Exec(`delete from tagacc where aid=?`, ID)
		if err != nil {
			panic(err)
		}

		stmt := `update acc set desc=? where id=?`
		_, err = DB.Exec(stmt, desc, ID)
		if err != nil {
			panic(err)
		}

		CreateTagInDesc(*desc, []int64{ID}, nil)
		CleanUpTags()
	}

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

	// 자식 계정들의 접두사(Prefix)도 함께 변경 (Cascade Rename)
	prefixOld := old + ":"
	prefixNeo := neo + ":"
	stmtChild := `update acc set name = ? || substr(name, ?) where name like ?`
	_, err = DB.Exec(stmtChild, prefixNeo, len([]rune(prefixOld))+1, prefixOld+"%")
	if err != nil {
		panic(err)
	}

	CleanUpTags()
	return ID
}
func DelAcc(name string) int64 {
	ChkDB()

	ID := GetAccByName(name)
	if ID == -1 {
		return -1
	}

	// 안전 장치: 거래 기록(Record)이 하나라도 존재하는 계정은 삭제 차단
	var count int
	err := DB.QueryRow(`select count(*) from record where aid=?`, ID).Scan(&count)
	if err != nil {
		panic(err)
	}
	if count > 0 {
		panic(fmt.Errorf("cannot delete account '%s' because it has %d associated transaction record(s)", name, count))
	}

	stmt := `delete from acc where id=?`
	_, err = DB.Exec(stmt, ID)
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
