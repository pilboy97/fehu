package core

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Rhymond/go-money"
)

func NewAcc(name string, desc string) (int64, error) {
	MustDB()

	stmt := `insert into acc(name, desc) values(?,?)`
	res, err := DB.Exec(stmt, name, desc)
	if err != nil {
		return 0, err
	}

	ret, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	CreateTagInDesc(desc, []int64{ret}, nil)
	return ret, nil
}
func GetAcc() ([]int64, error) {
	MustDB()

	stmt := `select id from acc order by id`
	rows, err := DB.Query(stmt)
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
func GetAccByID(aid int64) (Acc, error) {
	// TODO: ChkDB()는 panic을 발생시키므로, 여기서는 error를 반환하도록 변경해야 합니다.
	// 현재는 GetAccByID를 호출하는 곳에서 ChkDB()가 이미 호출되었다고 가정합니다.
	MustDB()

	var ret Acc
	stmt := `select id, name, desc from acc where id=?`
	err := DB.QueryRow(stmt, aid).Scan(
		&ret.ID,
		&ret.Name,
		&ret.Desc,
	)
	if err == sql.ErrNoRows {
		return Acc{ID: -1}, ErrCannotFind(fmt.Sprintf("account with ID %d", aid))
	}
	return ret, err
}
func GetAccByName(name string) (int64, error) {
	MustDB()

	var ret int64

	stmt := `select id from acc where name=?`
	err := DB.QueryRow(stmt, name).Scan(&ret)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrCannotFind(name)
		}
		return 0, err // Other database error
	}

	return ret, nil
}
func GetAccByPrefix(name string) ([]int64, error) {
	MustDB()

	rows, err := DB.Query(`select id from acc where name like ? order by id`, name+"%")
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
func GetAccByDesc(desc string) ([]int64, error) {
	MustDB()

	rows, err := DB.Query(`select id from acc where instr(desc,?) > 0 order by id`, desc)
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
func GetAccAmount(id int64) (*money.Money, error) {
	MustDB()

	acc, err := GetAccByID(id)
	if err != nil {
		// 계정을 찾을 수 없으면 0 금액 반환 (또는 에러 반환)
		return &money.Money{}, nil
	}
	if acc.ID == -1 { // GetAccByID가 -1을 반환하는 경우
		return &money.Money{}, nil
	}

	var ret *money.Money
	ret = money.New(0, Code)

	records, err := GetRecordByAID(id)
	if err != nil {
		return &money.Money{}, err
	}
	for _, rid := range records {
		record, err := GetRecordByID(rid)
		if err != nil {
			return &money.Money{}, err
		}
		ret, err = ret.Add(record.Amount)
		if err != nil {
			return &money.Money{}, err
		}
	}

	return ret, nil
}
func AltAcc(name string, desc *string) (int64, error) {
	MustDB()

	ID, err := GetAccByName(name)
	if err != nil {
		return 0, err
	}

	if desc == nil {
		return ID, nil
	}

	_, err = DB.Exec(`DELETE FROM tagacc WHERE aid=?`, ID)
	if err != nil {
		return 0, err
	}

	stmt := `update acc set desc=? where id=?`
	_, err = DB.Exec(stmt, desc, ID)
	if err != nil {
		return 0, err
	}

	CreateTagInDesc(*desc, []int64{ID}, nil)

	if err := CleanUpTags(); err != nil {
		return 0, err
	}

	return ID, nil
}
func AltRenameAcc(old, neo string) (int64, error) {
	MustDB()

	ID, err := GetAccByName(old)
	if err != nil {
		return 0, err
	}

	if ID2, err := GetAccByName(neo); err == nil && ID2 != -1 {
		return 0, ErrAlreadyExists("account: " + neo)
	}

	stmt := `update acc set name=? where id=?`
	_, err = DB.Exec(stmt, neo, ID)
	if err != nil {
		return 0, err
	}

	// 자식 계정들의 접두사(Prefix)도 함께 변경 (Cascade Rename)
	prefixOld := old + ":"
	prefixNeo := neo + ":"
	stmtChild := `update acc set name = ? || substr(name, ?) where name like ?`
	_, err = DB.Exec(stmtChild, prefixNeo, len(prefixOld)+1, prefixOld+"%") // len([]rune(prefixOld)) 대신 len(prefixOld) 사용
	if err != nil {
		return 0, err
	}

	if err := CleanUpTags(); err != nil {
		return 0, err
	}
	return ID, nil
}
func DelAcc(name string) (int64, error) {
	MustDB()

	ID, err := GetAccByName(name)
	if err != nil {
		return 0, err
	}

	// 안전 장치: 거래 기록(Record)이 하나라도 존재하는 계정은 삭제 차단
	var count int
	err = DB.QueryRow(`select count(*) from record where aid=?`, ID).Scan(&count)
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, fmt.Errorf("cannot delete account '%s' because it has %d associated transaction record(s)", name, count)
	}

	stmt := `delete from acc where id=?`
	_, err = DB.Exec(stmt, ID)
	if err != nil {
		return 0, err
	}

	return ID, nil
}

func PrintAccs(acc []int64) string {
	var ret = []string{}

	ret = append(ret, "      id|                    name|  amount|                    desc")
	for _, id := range acc {
		acc, err := GetAccByID(id)
		if err != nil {
			// 에러 처리: 계정을 찾을 수 없는 경우 스킵하거나 오류 메시지 포함
			continue
		}
		amount, err := GetAccAmount(id)
		if err != nil {
			panic(err)
		}

		ret = append(ret, fmt.Sprintf("%8d|%24s|%8s|%24s", acc.ID, acc.Name, amount.Absolute().Display(), acc.Desc))
	}
	ret = append(ret, fmt.Sprintf("%8d\tacc(s)\tfound", len(acc)))

	return strings.Join(ret, "\n")
}
