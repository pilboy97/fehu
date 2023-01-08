package core

import (
	"time"

	"github.com/Rhymond/go-money"
)

type Acc struct {
	ID   int64
	Name string
	Desc string
}
type Txn struct {
	ID     int64
	Desc   string
	Time   time.Time
	Record []Record
}
type Record struct {
	ID     int64
	TID    int64
	AID    int64
	Amount *money.Money
}
type Tag struct {
	ID   int64
	Name string
	Desc string
}
