package core

import "ast"

type Table struct {
	ids map[int64]struct{}
}

func NewTable(ids []int64) *Table {
	var set = make(map[int64]struct{})

	for _, id := range ids {
		set[id] = struct{}{}
	}

	return &Table{
		ids: set,
	}
}

func (t *Table) String() string {
	var ids = make([]int64, 0, len(t.ids))

	for id := range t.ids {
		ids = append(ids, id)
	}

	return PrintTxns(ids)
}
func (t *Table) Get() ast.Value {
	return t
}
func (t *Table) Count() int {
	return len(t.ids)
}
func (t *Table) Union(x *Table) *Table {
	set1, set2 := t.ids, x.ids
	set := make(map[int64]struct{})

	for id := range set1 {
		set[id] = struct{}{}
	}
	for id := range set2 {
		set[id] = struct{}{}
	}
	return &Table{ids: set}
}
func (t *Table) Intersect(x *Table) *Table {
	set1, set2 := t.ids, x.ids
	set := make(map[int64]struct{})

	for id := range set1 {
		if _, ok := set2[id]; ok {
			set[id] = struct{}{}
		}
	}
	return &Table{ids: set}
}
func (t *Table) Subtract(x *Table) *Table {
	set1, set2 := t.ids, x.ids
	set := make(map[int64]struct{})

	for id := range set1 {
		if _, ok := set2[id]; !ok {
			set[id] = struct{}{}
		}
	}
	return &Table{ids: set}
}
func (t *Table) XOR(x *Table) *Table {
	set1, set2 := t.ids, x.ids
	set := make(map[int64]struct{})

	for id := range set1 {
		if _, ok := set2[id]; !ok {
			set[id] = struct{}{}
		}
	}
	for id := range set2 {
		if _, ok := set1[id]; !ok {
			set[id] = struct{}{}
		}
	}
	return &Table{ids: set}
}
func (t *Table) Acc(ptn string) ast.List {
	var ret ast.List

	for id := range t.ids {
		r, err := GetRecordByTID(id)
		if err != nil {
			panic(err)
		}
		for _, rid := range r {
			record, err := GetRecordByID(rid)
			if err != nil {
				panic(err)
			}
			acc, err := GetAccByID(record.AID)
			if err != nil {
				panic(err)
			}

			if Search(acc.Name, ptn) {
				ret = ret.Append(
					ast.Num(record.Amount.AsMajorUnits()),
				)
			}
		}
	}

	return ret
}
func (t *Table) FilterPeriod(st, ed *int64) *Table {
	var ret = make([]int64, 0, len(t.ids))

	if st != nil && ed != nil {
		for tid := range t.ids {
			txn, err := GetTxnByID(tid)
			if err != nil {
				panic(err)
			}

			if *st < txn.Time && txn.Time < *ed {
				ret = append(ret, txn.ID)
			}
		}
	} else if st != nil {
		for tid := range t.ids {
			txn, err := GetTxnByID(tid)
			if err != nil {
				panic(err)
			}

			if *st < txn.Time {
				ret = append(ret, txn.ID)
			}
		}
	} else if ed != nil {
		for tid := range t.ids {
			txn, err := GetTxnByID(tid)
			if err != nil {
				panic(err)
			}

			if txn.Time < *ed {
				ret = append(ret, txn.ID)
			}
		}
	} else {
		for tid := range t.ids {
			ret = append(ret, tid)
		}
	}

	var set = make(map[int64]struct{})
	for _, tid := range ret {
		set[tid] = struct{}{}
	}
	return &Table{ids: set}
}
func (t *Table) ATag(name string) ast.List {
	var ret ast.List
	var tagid, err = GetTagByName(name)
	if err != nil {
		panic(err)
	}

	for id := range t.ids {
		r, err := GetRecordByTID(id)
		if err != nil {
			panic(err)
		}
		for _, rid := range r {
			record, err := GetRecordByID(rid)
			if err != nil {
				panic(err)
			}
			acc, err := GetAccByID(record.AID)
			if err != nil {
				panic(err)
			}

			if GetTagAcc(tagid, acc.ID) {
				ret = ret.Append(
					ast.Num(record.Amount.AsMajorUnits()),
				)
			}
		}
	}

	return ret
}
func (t *Table) TTag(name string) *Table {
	var ret = []int64{}
	var tagid, err = GetTagByName(name)
	if err != nil {
		panic(err)
	}

	for id := range t.ids {
		r, err := GetRecordByTID(id)
		if err != nil {
			panic(err)
		}
		for _, rid := range r {
			record, err := GetRecordByID(rid)
			if err != nil {
				panic(err)
			}

			if GetTagTxn(tagid, record.TID) {
				ret = append(ret, record.TID)
			}
		}
	}

	return NewTable(ret)
}

type Period struct {
	St, Ed *int64
}

func init() {
	Vars["__all__"] = ast.Computed{
		Fn: func() ast.Value {
			return NewTable(GetTxn())
		},
	}
}
