package ast

import (
	"fmt"
	"strings"
)

type Value interface {
	// 추상구문 트리의 노드의 값
	String() string
	// 노드를 문자열로 출력
	Get() Value
	// 노드의 실재 값을 반환
}

type Num float64

// 숫자 노드

func (n Num) String() string {
	//문자열로 출력
	return fmt.Sprint(float64(n))
}
func (n Num) Get() Value {
	return n
}

type Str string

func (s Str) String() string {
	return fmt.Sprint(string(s))
}
func (s Str) Get() Value {
	return s
}

type Bool bool

func (b Bool) String() string {
	return fmt.Sprint(bool(b))
}
func (b Bool) Get() Value {
	return b
}

type Void struct{}

func (v Void) String() string {
	return ""
}
func (v Void) Get() Value {
	return v
}

type Sym struct {
	Op int
}

func (s Sym) String() string {
	return fmt.Sprintf("<%s>", SymbolDesc[s.Op])
}
func (s Sym) Get() Value {
	return s
}

type Variable struct {
	Name string
}

func (v Variable) String() string {
	return fmt.Sprintf("<%s>", v.Name)
}
func (v Variable) Get() Value {
	return v
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
func (l List) Get() Value {
	return l
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

type Computed struct {
	Fn func() Value
}

func (c Computed) String() string {
	return c.Fn().String()
}
func (c Computed) Get() Value {
	return c.Fn()
}
