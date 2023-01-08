package core

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

var ErrDBIsNotOpened = errors.New("DB is not opened")
var DB *sql.DB

func ChkDB() {
	if len(DBPath) == 0 {
		panic(ErrDBIsNotOpened)
	}
}

func Open(path string) {
	DBPath = path

	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}

	_, err = DB.Exec(`pragma foreign_keys = on`)
	if err != nil {
		panic(err)
	}

	createAccStmt := `create table if not exists acc(
		id integer not null primary key autoincrement,
		name text not null unique,
		desc text
	)`

	_, err = DB.Exec(createAccStmt)
	if err != nil {
		panic(err)
	}

	createTxnStmt := `create table if not exists txn(
		id integer not null primary key autoincrement,
		desc text,
		time timestamp not null default CURRENT_TIMESTAMP
	)`

	_, err = DB.Exec(createTxnStmt)
	if err != nil {
		panic(err)
	}

	createRecordStmt := `create table if not exists record (
		id integer not null primary key autoincrement,
		tid integer not null,
		aid integer not null,
		amount integer not null,
		foreign key(tid) references txn(id) on delete cascade,
		foreign key(aid) references acc(id) on delete cascade
	)`

	_, err = DB.Exec(createRecordStmt)
	if err != nil {
		panic(err)
	}

	createTagStmt := `create table if not exists Tag(
		id integer not null primary key autoincrement,
		name text not null unique,
		desc text
	)`

	_, err = DB.Exec(createTagStmt)
	if err != nil {
		panic(err)
	}

	createTagAccStmt := `create table if not exists tagacc(
		tagid int64,
		aid int64,
		primary key(tagid, aid)
		foreign key(tagid) references tag(id) on delete cascade
		foreign key(aid) references acc(id) on delete cascade
	)`

	_, err = DB.Exec(createTagAccStmt)
	if err != nil {
		panic(err)
	}

	createTagTxnStmt := `create table if not exists tagtxn(
		tagid int64,
		tid int64,
		primary key(tagid, tid)
		foreign key(tagid) references tag(id) on delete cascade
		foreign key(tid) references txn(id) on delete cascade
	)`

	_, err = DB.Exec(createTagTxnStmt)
	if err != nil {
		panic(err)
	}

}
func Close() {
	DB.Close()
	DBPath = ""
}
