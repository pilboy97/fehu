package ast

const (
	NUM = iota
	STR
	TRUE
	FALSE
	NAME
	DEF
	POPEN
	PCLOSE
	ADD
	SUB
	MUL
	DIV
	NOT
	GR
	LE
	EQ
	NEQ
	GEQ
	LEQ
	AND
	OR
	FCALL
	COMMA
)

var SymbolDesc = map[int]string{
	NUM:    "NUMBER",
	STR:    "STRING",
	TRUE:   "TRUE",
	FALSE:  "FALSE",
	NAME:   "NAME",
	POPEN:  "(",
	PCLOSE: ")",
	ADD:    "+",
	SUB:    "-",
	MUL:    "*",
	DIV:    "/",
	NOT:    "!",
	GR:     ">",
	LE:     "<",
	EQ:     "==",
	NEQ:    "!=",
	GEQ:    ">=",
	LEQ:    "<=",
	DEF:    "def",
	AND:    "and",
	OR:     "or",
	FCALL:  "func call",
	COMMA:  ",",
}
var SymbolParam = map[int]int{
	ADD:   2,
	SUB:   2,
	MUL:   2,
	DIV:   2,
	EQ:    2,
	NEQ:   2,
	GR:    2,
	LE:    2,
	GEQ:   2,
	LEQ:   2,
	AND:   2,
	OR:    2,
	COMMA: 2,
	NOT:   1,
	FCALL: 2,
}
var SymbolPriority = map[int]int{
	ADD:   6,
	SUB:   6,
	MUL:   5,
	DIV:   5,
	EQ:    10,
	NEQ:   10,
	GR:    9,
	LE:    9,
	GEQ:   9,
	LEQ:   9,
	AND:   14,
	OR:    15,
	NOT:   3,
	FCALL: 2,
	COMMA: 17,
}
var Ops = map[string]int{
	"+":  ADD,
	"-":  SUB,
	"*":  MUL,
	"/":  DIV,
	"!":  NOT,
	">":  GR,
	"<":  LE,
	"==": EQ,
	"!=": NEQ,
	">=": GEQ,
	"<=": LEQ,
	",":  COMMA,
}
var RWords = map[string]int{
	"def":   DEF,
	"true":  TRUE,
	"false": FALSE,
	"and":   AND,
	"or":    OR,
}
