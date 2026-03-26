package core

import "ast"

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

func LoadAllDefsFromDB() {
	MustDB()

	var ret = make(map[string]string)
	stmt := `select name, stmt from vars`
	rows, err := DB.Query(stmt)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name, stmt string
		err = rows.Scan(&name, &stmt)
		if err != nil {
			panic(err)
		}

		ret[name] = stmt
	}

	for name, stmt := range ret {
		DefStmt(name, stmt)
	}
}
