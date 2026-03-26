package core

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

var ErrDBIsNotOpened = errors.New("DB is not opened")
var DB *sql.DB

func MustDB() {
	if DB == nil {
		if db := os.Getenv("FEHU_DB"); db != "" {
			if err := Open(db + ".db"); err != nil {
				panic(err)
			}
		}
	} else if err := DB.Ping(); err != nil {
		if db := os.Getenv("FEHU_DB"); db != "" {
			if err := Open(db + ".db"); err != nil {
				panic(err)
			}
		} else if DBPath != "" {
			if err := Open(DBPath); err != nil {
				panic(err)
			}
		}
	}

	if len(DBPath) == 0 {
		panic(ErrDBIsNotOpened)
	}
}

func Open(path string) error {
	DBPath = path

	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`pragma foreign_keys = on`)
	if err != nil {
		return err
	}

	createAccStmt := `create table if not exists acc(
		id integer not null primary key autoincrement,
		name text not null unique,
		desc text
	)`

	_, err = DB.Exec(createAccStmt)
	if err != nil {
		return err
	}

	createTxnStmt := `create table if not exists txn(
		id integer not null primary key autoincrement,
		"desc" text,
		"time" integer not null
	)`

	_, err = DB.Exec(createTxnStmt)
	if err != nil {
		return err
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
		return err
	}

	createTagStmt := `create table if not exists Tag(
		id integer not null primary key autoincrement,
		name text not null unique,
		desc text
	)`

	_, err = DB.Exec(createTagStmt)
	if err != nil {
		return err
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
		return err
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
		return err
	}

	createConfigStmt := `create table if not exists config(
		key text not null primary key,
		currency text not null
	)`

	_, err = DB.Exec(createConfigStmt)
	if err != nil {
		return err
	}

	createDefStmt := `create table if not exists vars (
		name text not null primary key,
		stmt text not null
	)`

	_, err = DB.Exec(createDefStmt)
	if err != nil {
		return err
	}

	if err := LoadConfig(); err != nil {
		return err
	}
	LoadAllDefsFromDB()

	return nil
}

func Close() {
	DB.Close()
	DBPath = ""
}
