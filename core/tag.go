package core

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
)

func FindTag(str string) []string {
	re := regexp.MustCompile(`#(\p{L}|\d|_)+`)
	return re.FindAllString(str, -1)
}
func CreateTagInDesc(desc string, aid []int64, tid []int64) {
	tags := FindTag(desc)

	for _, tag := range tags {
		tag = tag[1:]

		var id int64
		id = GetTagByName(tag)
		if id == -1 {
			id = NewTag(tag, "")
		}

		for _, a := range aid {
			NewTagAcc(id, a)
		}
		for _, t := range tid {
			NewTagTxn(id, t)
		}
	}
}

func NewTag(name string, desc string) int64 {
	ChkDB()

	stmt := `insert into tag(name, desc) values(?,?)`
	res, err := DB.Exec(stmt, name, desc)
	if err != nil {
		panic(err)
	}

	ret, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}
	return ret
}
func GetTag() []int64 {
	ChkDB()

	var ret []int64
	stmt := `select id from tag order by id`
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
func GetTagByID(aid int64) Tag {
	ChkDB()

	var ret Tag
	stmt := `select id, name, desc from tag where id=?`
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
func GetTagByName(name string) int64 {
	ChkDB()

	var ret int64

	stmt := `select id from tag where name=?`
	err := DB.QueryRow(stmt, name).Scan(&ret)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		return -1
	}

	return ret
}
func GetTagByDesc(desc string) []int64 {
	ChkDB()

	var ret []int64

	stmt := `select id from tag where instr(desc,?) > 0 order by id`
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
func GetTagAcc(tagid, aid int64) bool {
	stmt := `select count(*) from tagacc where tagid=? and aid=?`

	var count = 0
	if err := DB.QueryRow(stmt, tagid, aid).Scan(&count); err != nil {
		panic(err)
	}

	return count != 0
}
func GetTagTxn(tagid, tid int64) bool {
	stmt := `select count(*) from tagtxn where tagid=? and tid=?`

	var count = 0
	if err := DB.QueryRow(stmt, tagid, tid).Scan(&count); err != nil {
		panic(err)
	}

	return count != 0
}
func AltTag(name string, desc *string) int64 {
	ChkDB()

	ID := GetTagByName(name)
	if ID == -1 {
		return -1
	}

	if desc == nil {
		return ID
	}

	stmt := `update tag set desc=? where id=?`
	_, err := DB.Exec(stmt, desc, ID)
	if err != nil {
		panic(err)
	}

	return ID
}
func AltRenameTag(old, neo string) int64 {
	ChkDB()

	ID := GetTagByName(old)
	if ID == -1 {
		return -1
	}

	if ID2 := GetTagByName(neo); ID2 != -1 {
		return -2
	}

	stmt := `update tag set name=? where id=?`
	_, err := DB.Exec(stmt, neo, ID)
	if err != nil {
		panic(err)
	}

	// 기존 계정과 거래의 설명(desc) 원문에 적힌 해시태그 문자열도 일괄 치환
	// (예: #testtag2를 #testtag3으로 변경. #testtag2_extra 같은 부분 일치 방지)
	re := regexp.MustCompile(`(?i)#` + regexp.QuoteMeta(old) + `([^\p{L}\d_]|$)`)

	accRows, err := DB.Query(`select aid from tagacc where tagid=?`, ID)
	if err == nil {
		var aids []int64
		for accRows.Next() {
			var aid int64
			accRows.Scan(&aid)
			aids = append(aids, aid)
		}
		accRows.Close()

		for _, aid := range aids {
			acc := GetAccByID(aid)
			newDesc := re.ReplaceAllString(acc.Desc, "#"+neo+"${1}")
			DB.Exec(`update acc set desc=? where id=?`, newDesc, aid)
		}
	}

	txnRows, err := DB.Query(`select tid from tagtxn where tagid=?`, ID)
	if err == nil {
		var tids []int64
		for txnRows.Next() {
			var tid int64
			txnRows.Scan(&tid)
			tids = append(tids, tid)
		}
		txnRows.Close()

		for _, tid := range tids {
			txn := GetTxnByID(tid)
			newDesc := re.ReplaceAllString(txn.Desc, "#"+neo+"${1}")
			DB.Exec(`update txn set desc=? where id=?`, newDesc, tid)
		}
	}

	return ID
}
func CleanUpTags() {
	ChkDB()

	stmt := `delete from tag where id not in (select tagid from tagacc) and id not in (select tagid from tagtxn)`
	_, err := DB.Exec(stmt)
	if err != nil {
		panic(err)
	}
}
func DelTag(name string) int64 {
	ChkDB()

	ID := GetTagByName(name)
	if ID == -1 {
		return -1
	}

	stmt := `delete from tag where id=?`
	_, err := DB.Exec(stmt, ID)
	if err != nil {
		panic(err)
	}

	return ID
}
func NewTagAcc(tagid, aid int64) {
	stmt := `insert into tagacc(tagid, aid) values(?,?)`
	if _, err := DB.Exec(stmt, tagid, aid); err != nil {
		panic(err)
	}
}
func NewTagTxn(tagid, tid int64) {
	stmt := `insert into tagtxn(tagid, tid) values(?,?)`
	if _, err := DB.Exec(stmt, tagid, tid); err != nil {
		panic(err)
	}
}

func PrintTags(tag []int64) string {
	var ret = []string{}

	ret = append(ret, "      id|                    name|                    desc")
	for _, id := range tag {
		var tag = GetTagByID(id)

		ret = append(ret, fmt.Sprintf("%8d|%24s|%24s", tag.ID, tag.Name, tag.Desc))
	}
	ret = append(ret, fmt.Sprintf("%8d\ttag(s)\tfound", len(tag)))

	return strings.Join(ret, "\n")
}
