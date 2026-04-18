package core

import (
	"ast"
	"fmt"
)

var Vars = map[string]ast.Value{
	"count":     ast.Variable{Name: "count"},
	"sum":       ast.Variable{Name: "sum"},
	"avg":       ast.Variable{Name: "avg"},
	"max":       ast.Variable{Name: "max"},
	"min":       ast.Variable{Name: "min"},
	"acc":       ast.Variable{Name: "acc"},
	"between":   ast.Variable{Name: "between"},
	"union":     ast.Variable{Name: "union"},
	"intersect": ast.Variable{Name: "intersect"},
	"xor":       ast.Variable{Name: "xor"},
	"atag":      ast.Variable{Name: "atag"},
	"ttag":      ast.Variable{Name: "ttag"},
	"help":      ast.Variable{Name: "help"},
}

func CalcStmt(stmt string) ast.Value {
	tok := ast.Tokenize(stmt)
	ast := ast.NewAst(tok)

	return Calc(ast)
}

func DefStmt(name, stmt string) error {
	MustDB()

	tok := ast.Tokenize(stmt)
	ast := ast.NewAst(tok)

	if _, ok := Vars[name]; ok {
		return ErrAlreadyExists(name)
	}

	Vars[name] = Calc(ast)

	query := `insert into vars(name, stmt) values(?,?) on conflict(name) do update set stmt=excluded.stmt`
	_, err := DB.Exec(query, name, stmt)

	return err
}

func ChkVar(name string) error { // This function seems to check if a variable exists in the DB
	MustDB()

	stmt := `select name, stmt from vars where name=?`
	row := DB.QueryRow(stmt, name)
	var n, s string
	return row.Scan(&n, &s)
}

func LoadAllDefsFromDB() error {
	MustDB()

	rows, err := DB.Query(`select name, stmt from vars`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var defs = make(map[string]string)
	for rows.Next() {
		var name, stmt string
		if err = rows.Scan(&name, &stmt); err != nil {
			return err
		}
		defs[name] = stmt
	}

	for name, stmt := range defs {
		if _, ok := Vars[name]; ok {
			// 이미 메모리에 정의된 변수는 건너뜀 (같은 프로세스에서 DB를 재오픈하는 경우)
			continue
		}
		if err := DefStmt(name, stmt); err != nil {
			return fmt.Errorf("loading var %q: %w", name, err)
		}
	}
	return nil
}
