package ast

import (
	"errors"
	"fmt"
	"sort"
)

// 추상구문트리
type Ast struct {
	// 노드
	V Value
	// 문자열
	Str string
	// 부모 노드
	Par *Ast
	// 자식 노드
	Child []*Ast
}

// 피연산자 숫자가 맞지 않을 때 예외
var ErrOPUnmatched = errors.New("unmatched operator")

// 토큰으로 추상구문트리 생성
// 부모노드 리턴
func NewAst(stmt []Token) *Ast {
	// 연산자 수를 최대 토큰의 수만큼 있다 가정하고 미리 할당
	var ops = make([]Token, 0, len(stmt))
	// 토큰들을 순회
	for _, tok := range stmt {
		// 만약 연산자가 맞다면 토큰에 대한 연산자 우선순위가 존재함
		// 따라서 연산자 우선순위가 존재하는지 로 연산자 여부를 판단
		if _, ok := SymbolPriority[tok.Sym]; ok {
			//만약 연산자가 맞다면 배열에 추가
			ops = append(ops, tok)
		}
	}

	//연산자 배열을 우선순위에 따라 정렬
	sort.Slice(ops, func(i, j int) bool {
		if ops[i].Depth != ops[j].Depth {
			// 더 많은 괄호안에 있는 연산자가 우선순위가 높음
			return ops[i].Depth > ops[j].Depth
		} else if ops[i].Priority != ops[j].Priority {
			// 우선순위가 높은 연산자가 우선
			return ops[i].Priority < ops[j].Priority
		}

		// 먼저 오는 연산자가 우선
		return ops[i].Location < ops[j].Location
	})

	// 트리 노드수는 토큰 수 만큼 존재 함
	var node = make([]*Ast, len(stmt))
	//트리 초기화
	for i := 0; i < len(stmt); i++ {
		node[i] = &Ast{}
	}

	// 토큰들을 순회
	for i := range stmt {
		if stmt[i].Sym == POPEN {
			// 만약 여는 괄호라면
			// 바로 뒤의 노드가 부모
			node[i].Par = node[i+1]
		} else if stmt[i].Sym == PCLOSE {
			// 만약 닫는 괄호라면
			// 바로 앞의 노드가 부모
			node[i].Par = node[i-1]
		} else {
			//일반 토큰은 자신이 부모로 초기화
			node[i].Par = node[i]
		}

		node[i].Str = stmt[i].Str
	}

	for i := range stmt {
		// 추상구문트리의 노드를 초기화함
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
		//연산자 배열을 순회함
		//앞에서 우선순위로 정렬되었음

		// 연산자의 위치
		var i = op.Location
		// 노드 초기화
		node[i].V = Sym{Op: op.Sym}
		switch SymbolParam[op.Sym] {
		//연산자의 피연산자 수에 따라 분기
		case 2:
			//이항연산자라면
			var L, R int

			L = -1
			R = -1

			for l := i - 1; l >= 0; l-- {
				if node[l].Par == node[l] {
					// 바로 왼쪽 노드가 속한 트리의 루트가 노드의 왼쪽 자식
					L = l
					break
				}
			}
			for r := i + 1; r < len(stmt); r++ {
				if node[r].Par == node[r] {
					// 바로 오른쪽 노드가 속한 트리의 루트가 노드의 오른쪽 자식
					R = r
					break
				}
			}

			if L == -1 || R == -1 {
				// 왼쪽과 오른쪽 자식 노드를 찾지 못했다면
				if op.Sym == FCALL && R == -1 {
					// 만약 연산자가 함수 호출이라면
					node[L].Par = node[i]
					node[i].Child = []*Ast{node[L]}

				} else {
					//아니라면 예외 처리
					panic(ErrOPUnmatched)
				}
			} else {
				node[L].Par = node[i]
				node[R].Par = node[i]

				node[i].Child = []*Ast{node[L], node[R]}
			}
		case 1:
			// 단항 연산자라면
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
