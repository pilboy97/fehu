package core

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

var ErrDBIsNotOpened = errors.New("DB is not opened")
var DB *sql.DB

func ChkDB() {
	if DB == nil {
		if db := os.Getenv("FEHU_DB"); db != "" {
			Open(db + ".db")
		}
	} else if err := DB.Ping(); err != nil {
		if db := os.Getenv("FEHU_DB"); db != "" {
			Open(db + ".db")
		} else if DBPath != "" {
			Open(DBPath)
		}
	}

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
		"desc" text,
		"time" integer not null
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
		tagid integer,
		aid integer,
		primary key(tagid, aid),
		foreign key(tagid) references tag(id) on delete cascade,
		foreign key(aid) references acc(id) on delete cascade
	)`

	_, err = DB.Exec(createTagAccStmt)
	if err != nil {
		panic(err)
	}

	createTagTxnStmt := `create table if not exists tagtxn(
		tagid integer,
		tid integer,
		primary key(tagid, tid),
		foreign key(tagid) references tag(id) on delete cascade,
		foreign key(tid) references txn(id) on delete cascade
	)`

	_, err = DB.Exec(createTagTxnStmt)
	if err != nil {
		panic(err)
	}

	createConfigStmt := `create table if not exists config(
		key text not null primary key,
		currency text not null
	)`

	_, err = DB.Exec(createConfigStmt)
	if err != nil {
		panic(err)
	}

	InitConfig()
}
func Close() {
	DB.Close()
	DBPath = ""
}
