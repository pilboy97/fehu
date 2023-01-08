package core

import (
	"ast"
	"errors"
	"math"
	"strings"
)

var ErrUnknownVariable = errors.New("unknown variable")
var ErrWrongType = errors.New("wrong type")
var ErrDividZero = errors.New("divid by zero")
var ErrNotMatchFunc = errors.New("cannot match function parameter")

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
}

func CalcStmt(stmt string) ast.Value {
	tok := ast.Tokenize(stmt)
	ast := ast.NewAst(tok)

	return Calc(ast)
}
func DefStmt(name, stmt string) {
	tok := ast.Tokenize(stmt)
	ast := ast.NewAst(tok)

	if _, ok := Vars[name]; ok {
		panic(ErrAlreadyExists)
	}

	Vars[name] = Calc(ast)
}
func Calc(root *ast.Ast) ast.Value {
	switch v := root.V.(type) {
	case ast.Num:
		return v
	case ast.Str:
		return v
	case ast.Bool:
		return v
	case ast.Variable:
		if v, ok := Vars[v.Name]; ok {
			return v
		} else {
			panic(ErrUnknownVariable)
		}
	case ast.Sym:
		switch v.Op {
		case ast.ADD:
			L, R := Calc(root.Child[0]), Calc(root.Child[1])
			if l, ok := L.(ast.Num); ok {
				if r, ok := R.(ast.Num); ok {
					return ast.Num(float64(l) + float64(r))
				}
			}
			if l, ok := L.(ast.Str); ok {
				if r, ok := R.(ast.Str); ok {
					return ast.Str(strings.Join([]string{string(l), string(r)}, ""))
				}
			}

		case ast.SUB:
			L, R := Calc(root.Child[0]), Calc(root.Child[1])
			if l, ok := L.(ast.Num); ok {
				if r, ok := R.(ast.Num); ok {
					return ast.Num(float64(l) - float64(r))
				}
			}

		case ast.MUL:
			L, R := Calc(root.Child[0]), Calc(root.Child[1])
			if l, ok := L.(ast.Num); ok {
				if r, ok := R.(ast.Num); ok {
					return ast.Num(float64(l) * float64(r))
				}
			}

		case ast.DIV:
			L, R := Calc(root.Child[0]), Calc(root.Child[1])
			if l, ok := L.(ast.Num); ok {
				if r, ok := R.(ast.Num); ok {
					if float64(r) == 0 {
						panic(ErrDividZero)
					}
					return ast.Num(float64(l) / float64(r))
				}
			}

		case ast.EQ:
			L, R := Calc(root.Child[0]), Calc(root.Child[1])
			if l, ok := L.(ast.Num); ok {
				if r, ok := R.(ast.Num); ok {
					return ast.Bool(float64(l) == float64(r))
				}
			}

		case ast.NEQ:
			L, R := Calc(root.Child[0]), Calc(root.Child[1])
			if l, ok := L.(ast.Num); ok {
				if r, ok := R.(ast.Num); ok {
					return ast.Bool(float64(l) != float64(r))
				}
			}

		case ast.GR:
			L, R := Calc(root.Child[0]), Calc(root.Child[1])
			if l, ok := L.(ast.Num); ok {
				if r, ok := R.(ast.Num); ok {
					return ast.Bool(float64(l) > float64(r))
				}
			}

		case ast.LE:
			L, R := Calc(root.Child[0]), Calc(root.Child[1])
			if l, ok := L.(ast.Num); ok {
				if r, ok := R.(ast.Num); ok {
					return ast.Bool(float64(l) < float64(r))
				}
			}

		case ast.GEQ:
			L, R := Calc(root.Child[0]), Calc(root.Child[1])
			if l, ok := L.(ast.Num); ok {
				if r, ok := R.(ast.Num); ok {
					return ast.Bool(float64(l) >= float64(r))
				}
			}

		case ast.LEQ:
			L, R := Calc(root.Child[0]), Calc(root.Child[1])
			if l, ok := L.(ast.Num); ok {
				if r, ok := R.(ast.Num); ok {
					return ast.Bool(float64(l) <= float64(r))
				}
			}

		case ast.AND:
			L, R := Calc(root.Child[0]), Calc(root.Child[1])
			if l, ok := L.(ast.Bool); ok {
				if r, ok := R.(ast.Bool); ok {
					return ast.Bool(bool(l) && bool(r))
				}
			}

		case ast.OR:
			L, R := Calc(root.Child[0]), Calc(root.Child[1])
			if l, ok := L.(ast.Bool); ok {
				if r, ok := R.(ast.Bool); ok {
					return ast.Bool(bool(l) || bool(r))
				}
			}

		case ast.COMMA:
			L, R := Calc(root.Child[0]), Calc(root.Child[1])
			ret := ast.List{}

			return ret.Append(L).Append(R)

		case ast.FCALL:
			return FCall(root)

		case ast.NOT:
			R := Calc(root.Child[0])
			if r, ok := R.(ast.Bool); ok {
				return ast.Bool(!bool(r))
			}
		}
	default:
	}
	panic(ErrWrongType)
}
func FCall(root *ast.Ast) ast.Value {
	var L, R ast.Value
	L = Calc(root.Child[0])
	if len(root.Child) >= 2 {
		R = Calc(root.Child[1])
	} else {
		R = ast.List{}
	}

	if _, ok := L.(ast.Variable); !ok {
		panic(ErrWrongType)
	}
	if _, ok := R.(ast.List); !ok {
		R = ast.List{}.Append(R)
	}

	l := L.(ast.Variable)
	r := R.(ast.List)

	if _, ok := Vars[l.Name]; !ok {
		panic(ErrUnknownVariable)
	}

	switch l.Name {
	case "sum":
		if len(r.List()) == 0 {
			return ast.Num(0)
		}
		var res float64 = 0
		for _, elem := range r.List() {
			if v, ok := elem.(ast.Num); ok {
				res += float64(v)
			} else {
				panic(ErrWrongType)
			}
		}

		return ast.Num(res)
	case "avg":
		if len(r.List()) == 0 {
			return ast.Num(math.NaN())
		}
		var res float64 = 0
		for _, elem := range r.List() {
			if v, ok := elem.(ast.Num); ok {
				res += float64(v)
			} else {
				panic(ErrWrongType)
			}
		}

		return ast.Num(res / float64(len(r.List())))
	case "count":
		if len(r.List()) == 0 {
			return ast.Num(0)
		}
		if t, ok := r.List()[0].(*Table); ok {
			return ast.Num(t.Count())
		}
		return ast.Num(float64(len(r.List())))
	case "max":
		if len(r.List()) == 0 {
			return ast.Num(math.NaN())
		}
		var res float64 = float64(r.List()[0].(ast.Num))
		for _, elem := range r.List() {
			if v, ok := elem.(ast.Num); ok {
				res = math.Max(res, float64(v))
			} else {
				panic(ErrWrongType)
			}
		}

		return ast.Num(res)
	case "min":
		if len(r.List()) == 0 {
			return ast.Num(math.NaN())
		}
		var res float64 = float64(r.List()[0].(ast.Num))
		for _, elem := range r.List() {
			if v, ok := elem.(ast.Num); ok {
				res = math.Min(res, float64(v))
			} else {
				panic(ErrWrongType)
			}
		}

		return ast.Num(res)
	case "acc":
		if len(r.List()) != 2 {
			panic(ErrNotMatchFunc)
		}

		var table *Table
		var name string

		if v, ok := r.List()[0].(*Table); ok {
			table = v
		} else {
			panic(ErrWrongType)
		}
		if v, ok := r.List()[1].(ast.Str); ok {
			name = string(v)
		} else {
			panic(ErrWrongType)
		}

		return table.Acc(name)
	case "atag":
		if len(r.List()) != 2 {
			panic(ErrNotMatchFunc)
		}

		var table *Table
		var name string

		if v, ok := r.List()[0].(*Table); ok {
			table = v
		} else {
			panic(ErrWrongType)
		}
		if v, ok := r.List()[1].(ast.Str); ok {
			name = string(v)
		} else {
			panic(ErrWrongType)
		}

		return table.ATag(name)
	case "ttag":
		if len(r.List()) != 2 {
			panic(ErrNotMatchFunc)
		}

		var table *Table
		var name string

		if v, ok := r.List()[0].(*Table); ok {
			table = v
		} else {
			panic(ErrWrongType)
		}
		if v, ok := r.List()[1].(ast.Str); ok {
			name = string(v)
		} else {
			panic(ErrWrongType)
		}

		return table.TTag(name)
	case "between":
		if len(r.List()) != 3 {
			panic(ErrNotMatchFunc)
		}

		var table *Table
		var st, ed string

		if v, ok := r.List()[0].(*Table); ok {
			table = v
		} else {
			panic(ErrWrongType)
		}
		if v, ok := r.List()[1].(ast.Str); ok {
			st = string(v)
		} else {
			panic(ErrWrongType)
		}
		if v, ok := r.List()[2].(ast.Str); ok {
			ed = string(v)
		} else {
			panic(ErrWrongType)
		}

		P := ParsePeriod(st, ed)

		return table.FilterPeriod(P.St, P.Ed)
	case "union":
		list := r.List()
		if len(list) < 2 {
			panic(ErrNotMatchFunc)
		}

		var ret *Table = nil

		if v1, ok := list[0].(*Table); ok {
			if v2, ok := list[1].(*Table); ok {
				ret = v1.Union(v2)
			}
		}

		if ret == nil {
			panic(ErrWrongType)
		}

		for i := 2; i < len(list); i++ {
			if v, ok := list[i].(*Table); ok {
				ret = ret.Union(v)
			} else {
				panic(ErrWrongType)
			}
		}

		return ret
	case "intersect":
		list := r.List()
		if len(list) < 2 {
			panic(ErrNotMatchFunc)
		}

		var ret *Table = nil

		if v1, ok := list[0].(*Table); ok {
			if v2, ok := list[1].(*Table); ok {
				ret = v1.Intersect(v2)
			}
		}

		if ret == nil {
			panic(ErrWrongType)
		}

		for i := 2; i < len(list); i++ {
			if v, ok := list[i].(*Table); ok {
				ret = ret.Intersect(v)
			} else {
				panic(ErrWrongType)
			}
		}

		return ret
	case "xor":
		list := r.List()
		if len(list) < 2 {
			panic(ErrNotMatchFunc)
		}

		var ret *Table = nil

		if v1, ok := list[0].(*Table); ok {
			if v2, ok := list[1].(*Table); ok {
				ret = v1.XOR(v2)
			}
		}

		if ret == nil {
			panic(ErrWrongType)
		}

		for i := 2; i < len(list); i++ {
			if v, ok := list[i].(*Table); ok {
				ret = ret.XOR(v)
			} else {
				panic(ErrWrongType)
			}
		}

		return ret
	default:
		panic(ErrWrongName)
	}
}
