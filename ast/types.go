package ast

import (
	"fmt"
	"strings"
)

type Value interface {
	String() string
}

type Num float64

func (n Num) String() string {
	return fmt.Sprint(float64(n))
}

type Str string

func (s Str) String() string {
	return fmt.Sprint(string(s))
}

type Bool bool

func (b Bool) String() string {
	return fmt.Sprint(bool(b))
}

type Void struct{}

func (v Void) String() string {
	return ""
}

type Sym struct {
	Op int
}

func (s Sym) String() string {
	return fmt.Sprintf("<%s>", SymbolDesc[s.Op])
}

type Variable struct {
	Name string
}

func (v Variable) String() string {
	return fmt.Sprintf("<%s>", v.Name)
}

type List struct {
	list []Value
}

func (l List) List() []Value {
	return l.list
}
func (l List) String() string {
	ret := []string{"["}
	if len(l.list) > 0 {
		for i := 0; i < len(l.list)-1; i++ {
			ret = append(ret, fmt.Sprintf("%s, ", l.list[i].String()))
		}
		ret = append(ret, l.list[len(l.list)-1].String())
	}
	ret = append(ret, "]")

	return strings.Join(ret, "")
}
func (l List) Append(v Value) List {
	ret := make([]Value, 0)
	ret = append(ret, l.list...)

	if l2, ok := v.(List); ok {
		ret = append(ret, l2.list...)
		return List{list: ret}
	} else {
		ret = append(ret, v)
	}
	return List{list: ret}
}
