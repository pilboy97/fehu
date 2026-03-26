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
		id, err := GetTagByName(tag)
		if err != nil {
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
	MustDB()

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
	MustDB()

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
func GetTagByID(aid int64) (Tag, error) {
	// TODO: ChkDB()는 panic을 발생시키므로, 여기서는 error를 반환하도록 변경해야 합니다.
	// 현재는 GetTagByID를 호출하는 곳에서 ChkDB()가 이미 호출되었다고 가정합니다.
	MustDB()

	var ret Tag
	stmt := `select id, name, desc from tag where id=?`
	err := DB.QueryRow(stmt, aid).Scan(
		&ret.ID,
		&ret.Name,
		&ret.Desc,
	)
	if err == sql.ErrNoRows {
		return Tag{ID: -1}, ErrCannotFind(fmt.Sprintf("tag with ID %d", aid))
	}
	return ret, err
}
func GetTagByName(name string) (int64, error) {
	MustDB()

	var ret int64

	stmt := `select id from tag where name=?`
	err := DB.QueryRow(stmt, name).Scan(&ret)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		return 0, ErrCannotFind(name)
	}

	return ret, nil
}
func GetTagByDesc(desc string) []int64 {
	MustDB()

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
		return false // Consider returning error here as well
	}

	return count != 0
}
func GetTagTxn(tagid, tid int64) bool {
	// TODO: ChkDB()는 panic을 발생시키므로, 여기서는 error를 반환하도록 변경해야 합니다.
	// 현재는 GetTagTxn을 호출하는 곳에서 ChkDB()가 이미 호출되었다고 가정합니다.
	MustDB()

	stmt := `select count(*) from tagtxn where tagid=? and tid=?`

	var count = 0
	if err := DB.QueryRow(stmt, tagid, tid).Scan(&count); err != nil {
		return false // Consider returning error here as well
	}

	return count != 0 // TODO: Should this return an error if DB.QueryRow fails?
}
func AltTag(name string, desc *string) (int64, error) {
	MustDB()

	ID, err := GetTagByName(name)
	if err != nil {
		return 0, err
	}

	if desc == nil {
		return ID, nil
	}

	stmt := `update tag set desc=? where id=?`
	_, err = DB.Exec(stmt, desc, ID)
	if err != nil {
		return 0, err
	}

	return ID, nil
}
func AltRenameTag(old, neo string) (int64, error) {
	MustDB()

	ID, err := GetTagByName(old)
	if err != nil {
		return 0, err
	}

	if _, err := GetTagByName(neo); err == nil { // If new tag name already exists
		return 0, ErrAlreadyExists("tag: " + neo)
	}

	stmt := `update tag set name=? where id=?`
	_, err = DB.Exec(stmt, neo, ID)
	if err != nil {
		return 0, err
	}

	// 기존 계정과 거래의 설명(desc) 원문에 적힌 해시태그 문자열도 일괄 치환
	// (예: #testtag2를 #testtag3으로 변경. #testtag2_extra 같은 부분 일치 방지)
	re := regexp.MustCompile(`(?i)#` + regexp.QuoteMeta(old) + `([^\p{L}\d_]|$)`)

	accRows, err := DB.Query(`select aid from tagacc where tagid=?`, ID)
	if err != nil {
		return 0, err
	}
	defer accRows.Close()

	var aids []int64
	for accRows.Next() {
		var aid int64
		if err := accRows.Scan(&aid); err != nil {
			return 0, err
		}
		aids = append(aids, aid)
	}

	if len(aids) > 0 {
		for _, aid := range aids {
			acc, err := GetAccByID(aid)
			if err != nil {
				return 0, err
			}

			newDesc := re.ReplaceAllString(acc.Desc, "#"+neo+"${1}")
			DB.Exec(`update acc set desc=? where id=?`, newDesc, aid)
		}
	}

	txnRows, err := DB.Query(`SELECT tid FROM tagtxn WHERE tagid=?`, ID)
	if err != nil {
		return 0, err
	}
	defer txnRows.Close()

	var tids []int64
	for txnRows.Next() {
		var tid int64
		if err := txnRows.Scan(&tid); err != nil {
			return 0, err
		}
		tids = append(tids, tid)
	}
	if len(tids) > 0 {

		for _, tid := range tids {
			txn, err := GetTxnByID(tid)
			if err != nil {
				return 0, err
			}

			newDesc := re.ReplaceAllString(txn.Desc, "#"+neo+"${1}")
			DB.Exec(`update txn set desc=? where id=?`, newDesc, tid)
		}
	}

	return ID, nil
}
func CleanUpTags() error {
	MustDB()

	stmt := `delete from tag where id not in (select tagid from tagacc) and id not in (select tagid from tagtxn)`
	_, err := DB.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}
func DelTag(name string) (int64, error) {
	MustDB()

	ID, err := GetTagByName(name)
	if err != nil {
		return 0, err
	}

	stmt := `delete from tag where id=?`
	_, err = DB.Exec(stmt, ID)
	if err != nil {
		return 0, err
	}

	return ID, nil
}
func NewTagAcc(tagid, aid int64) error {
	MustDB()

	stmt := `insert into tagacc(tagid, aid) values(?,?)`
	if _, err := DB.Exec(stmt, tagid, aid); err != nil {
		return err
	}

	return nil
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
		tag, err := GetTagByID(id)
		if err != nil {
			// 에러 처리: 태그를 찾을 수 없는 경우 스킵하거나 오류 메시지 포함
			continue
		}

		ret = append(ret, fmt.Sprintf("%8d|%24s|%24s", tag.ID, tag.Name, tag.Desc))
	}
	ret = append(ret, fmt.Sprintf("%8d\ttag(s)\tfound", len(tag)))

	return strings.Join(ret, "\n")
}
