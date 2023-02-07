package ast

import (
	"errors"
	"strconv"
	"unicode"
)

type Token struct {
	Sym      int
	Depth    int
	Priority int
	Location int
	Str      string
	Param    int
	Param2   float64
}

var ErrPairUnmatched = errors.New("unmatched pair")
var ErrUnknownSymbol = errors.New("unknown symbol")
var ErrQuoteUnmatched = errors.New("unmatched quote")

func inRange(x rune, st, ed rune) bool {
	if st > ed {
		st, ed = ed, st
	}
	return st <= x && x <= ed
}
func Tokenize(str string) []Token {
	runes := []rune(str)

	stack := []int{}
	ret := make([]Token, 0, len(runes))

	for i := 0; i < len(runes); i++ {
		if unicode.IsSpace(runes[i]) {
			continue
		}

		var tok Token

		switch {
		case runes[i] == '\'':
			var done = false
			for j := i + 1; j < len(runes); j++ {
				if runes[j] == '\'' {
					str = string(runes[i+1 : j])

					tok.Sym = STR
					tok.Str = str
					tok.Location = len(ret)
					tok.Depth = len(stack)

					i = j
					done = true

					break
				}
			}
			if !done {
				panic(ErrQuoteUnmatched)
			}
		case inRange(runes[i], '0', '9'):
			var done = false
			for j := i + 1; j < len(runes); j++ {
				if !(inRange(runes[j], '0', '9') || runes[j] == '.') {
					str = string(runes[i:j])
					val, err := strconv.ParseFloat(str, 64)
					if err != nil {
						panic(err)
					}

					tok.Sym = NUM
					tok.Str = str
					tok.Location = len(ret)
					tok.Depth = len(stack)
					tok.Param2 = val

					i = j - 1
					done = true
					break
				}
			}
			if !done {
				str = string(runes[i:])
				val, err := strconv.ParseFloat(str, 64)
				if err != nil {
					panic(err)
				}

				tok.Sym = NUM
				tok.Str = string(runes[i:])
				tok.Location = len(ret)
				tok.Depth = len(stack)
				tok.Param2 = val

				i = len(runes)
			}
		case runes[i] == '_' || inRange(runes[i], 'a', 'z') || inRange(runes[i], 'A', 'Z'):
			var done = false
			for j := i + 1; j < len(runes); j++ {
				if !(runes[j] == '_' ||
					inRange(runes[j], '0', '9') ||
					inRange(runes[j], 'a', 'z') ||
					inRange(runes[j], 'A', 'Z')) {
					tok.Sym = NAME
					tok.Str = string(runes[i:j])
					tok.Location = len(ret)
					tok.Depth = len(stack)

					i = j - 1
					done = true
					break
				}
			}
			if !done {
				tok.Sym = NAME
				tok.Str = string(runes[i:])
				tok.Location = i
				tok.Depth = len(stack)
				i = len(runes)
			}
		case runes[i] == '(':
			tok.Sym = POPEN
			tok.Str = "("
			tok.Location = len(ret)
			tok.Depth = len(stack)

			stack = append(stack, i)
		case runes[i] == ')':
			if len(stack) <= 0 {
				panic(ErrPairUnmatched)
			}
			pair := stack[len(stack)-1]

			tok.Sym = PCLOSE
			tok.Str = ")"
			tok.Location = len(ret)
			tok.Depth = len(stack)
			tok.Param = pair

			stack = stack[:len(stack)-1]
		default:
			done := false
			if i+1 < len(runes) {
				str := string(runes[i : i+2])
				if _, ok := Ops[str]; ok {
					done = true
					tok.Sym = Ops[str]
					tok.Str = str
					tok.Location = len(ret)
					tok.Priority = SymbolPriority[Ops[str]]
					tok.Depth = len(stack)

					i = i + 1
				}
			}
			if !done {
				str := string(runes[i])
				if _, ok := Ops[str]; !ok {
					panic(ErrUnknownSymbol)
				}

				tok.Sym = Ops[str]
				tok.Str = str
				tok.Location = len(ret)
				tok.Priority = SymbolPriority[Ops[str]]
				tok.Depth = len(stack)
			}
		}

		ret = append(ret, tok)
	}

	if len(stack) > 0 {
		panic(ErrPairUnmatched)
	}

	for i := 0; i < len(ret); i++ {
		if ret[i].Sym == NAME {
			if v, ok := RWords[ret[i].Str]; ok {
				ret[i].Sym = v
			}
		} else if i > 0 && ret[i].Sym == POPEN {
			if ret[i-1].Sym == NAME {
				ret[i].Str = "Func call"
				ret[i].Sym = FCALL
			}
		}
	}

	return ret
}
