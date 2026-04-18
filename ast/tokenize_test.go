package ast

import (
	"testing"
)

func TestTokenizeLocation(t *testing.T) {
	// 모든 토큰의 Location은 토큰 배열 인덱스(0, 1, 2...)여야 한다.
	// 버그: NAME 토큰이 문자열 끝에서 끝날 때 rune 위치(i)를 썼던 문제.
	cases := []struct {
		input string
		locs  []int // 각 토큰의 기대 Location
	}{
		{"abc", []int{0}},           // NAME만, EOF에서 끝남
		{"abc def", []int{0, 1}},    // NAME NAME
		{"abc+def", []int{0, 1, 2}}, // NAME OP NAME(EOF에서 끝남)
		{"1+abc", []int{0, 1, 2}},   // NUM OP NAME(EOF에서 끝남)
		{"abc+1", []int{0, 1, 2}},   // NAME OP NUM
	}

	for _, c := range cases {
		tokens := Tokenize(c.input)
		if len(tokens) != len(c.locs) {
			t.Errorf("input %q: got %d tokens, want %d", c.input, len(tokens), len(c.locs))
			continue
		}
		for i, tok := range tokens {
			if tok.Location != c.locs[i] {
				t.Errorf("input %q token[%d] (%q): Location=%d, want %d",
					c.input, i, tok.Str, tok.Location, c.locs[i])
			}
		}
	}
}

func TestTokenizeLocationConsistency(t *testing.T) {
	// Location 값이 항상 0부터 연속된 정수여야 한다 (rune 위치가 아님)
	inputs := []string{
		"sum acc",
		"count abc",
		"x+y+z",
		"foo",
	}

	for _, input := range inputs {
		tokens := Tokenize(input)
		for i, tok := range tokens {
			if tok.Location != i {
				t.Errorf("input %q: token[%d] %q has Location=%d, want %d",
					input, i, tok.Str, tok.Location, i)
			}
		}
	}
}
