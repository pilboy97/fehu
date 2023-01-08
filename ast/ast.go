package ast

import (
	"errors"
	"fmt"
	"sort"
)

type Ast struct {
	V     Value
	Str   string
	Par   *Ast
	Child []*Ast
}

var ErrOPUnmatched = errors.New("unmatched operator")

func NewAst(stmt []Token) *Ast {
	var ops = make([]Token, 0, len(stmt))
	for _, tok := range stmt {
		if _, ok := SymbolPriority[tok.Sym]; ok {
			ops = append(ops, tok)
		}
	}

	sort.Slice(ops, func(i, j int) bool {
		if ops[i].Depth != ops[j].Depth {
			return ops[i].Depth > ops[j].Depth
		} else if ops[i].Priority != ops[j].Priority {
			return ops[i].Priority < ops[j].Priority
		}

		return ops[i].Location < ops[j].Location
	})

	var node = make([]*Ast, len(stmt))
	for i := 0; i < len(stmt); i++ {
		node[i] = &Ast{}
	}

	for i := range stmt {
		if stmt[i].Sym == POPEN {
			node[i].Par = node[i+1]
		} else if stmt[i].Sym == PCLOSE {
			node[i].Par = node[i-1]
		} else {
			node[i].Par = node[i]
		}

		node[i].Str = stmt[i].Str
	}

	for i := range stmt {
		switch stmt[i].Sym {
		case NUM:
			node[i].V = Num(stmt[i].Param2)
		case STR:
			node[i].V = Str(stmt[i].Str)
		case TRUE:
			node[i].V = Bool(true)
		case FALSE:
			node[i].V = Bool(false)
		case NAME:
			node[i].V = Variable{Name: stmt[i].Str}
		default:
			node[i].V = Void{}
		}
	}
	for _, op := range ops {
		var i = op.Location
		node[i].V = Sym{Op: op.Sym}
		switch SymbolParam[op.Sym] {
		case 2:
			var L, R int

			L = -1
			R = -1

			for l := i - 1; l >= 0; l-- {
				if node[l].Par == node[l] {
					L = l
					break
				}
			}
			for r := i + 1; r < len(stmt); r++ {
				if node[r].Par == node[r] {
					R = r
					break
				}
			}

			if L == -1 || R == -1 {
				if op.Sym == FCALL && R == -1 {
					node[L].Par = node[i]
					node[i].Child = []*Ast{node[L]}

				} else {
					panic(ErrOPUnmatched)
				}
			} else {
				node[L].Par = node[i]
				node[R].Par = node[i]

				node[i].Child = []*Ast{node[L], node[R]}
			}
		case 1:
			var R int = -1
			for R = i + 1; R < len(stmt); R++ {
				if node[R].Par == node[R] {
					break
				}
			}
			if R == -1 {
				panic(ErrOPUnmatched)
			}

			node[R].Par = node[i]
			node[i].Child = []*Ast{node[R]}
		default:
		}
	}

	var root = node[0]
	for root.Par != root {
		root = root.Par
	}

	return root
}
func (root *Ast) Tour() {
	for _, ch := range root.Child {
		ch.Tour()
	}
	fmt.Println(root.Str)
}
