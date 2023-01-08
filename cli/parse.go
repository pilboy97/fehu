package cli

import (
	"fmt"
	"strings"
)

type UnknownStateError struct {
	What string
}

func (e *UnknownStateError) Error() string {
	return fmt.Sprintf("cannot parse %s", e.What)
}

type UnknownFlagError struct {
	S string
	F string
}

func (e *UnknownFlagError) Error() string {
	return fmt.Sprintf("%s : cannot parse flag : %s", e.S, e.F)
}

type FlagVar struct {
	F *Flag
	V string
}

type Cmd struct {
	St *State
	Fl []FlagVar
	Pa []string
}

type Parser struct {
	root *State
}

func NewParser(root *State) *Parser {
	return &Parser{
		root: root,
	}
}
func (p *Parser) Parse(str string) (Cmd, error) {
	if len(str) == 0 {
		return Cmd{St: p.root}, nil
	}

	token := TokenizeCommand(str)

	stack := []*State{p.root}
	result := Cmd{
		St: p.root,
		Fl: []FlagVar{},
		Pa: []string{},
	}
	set := make(map[*State]struct{})

	for i := 0; i < len(token); i++ {
		if token[i][0] == '"' && token[i][len(token[i])-1] == '"' {
			token[i] = token[i][1 : len(token[i])-1]
		}
	}

	var i = 0
	for len(stack) > 0 && i < len(token) {

		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		set[top] = struct{}{}

		for _, child := range top.next {
			if _, exists := set[child]; !exists {
				set[child] = struct{}{}

				ok, err := child.Match(token[i])
				if err != nil {
					return Cmd{}, err
				}

				if ok {
					stack = append(stack, child)
					result.St = child
					i++

					break
				}
			}
		}

	}

	if !result.St.Action && i < len(token) && !strings.HasPrefix(token[i], "-") {
		return Cmd{}, &UnknownStateError{What: token[i]}
	}

	for ; i < len(token); i++ {

		if !strings.HasPrefix(token[i], "-") {
			break
		}

		var found = false
		cur := token[i]
		values := ParseFlag(cur)

		for _, f := range result.St.Flags {

			res, err := f.Match(cur)

			if err != nil {
				return Cmd{}, err
			}

			if res {
				found = true

				var v string

				if len(values) == 2 {
					v = values[1]
				} else {
					v = "true"
				}

				if len(v) >= 2 && v[0] == '"' && v[len(v)-1] == '"' {
					v = v[1 : len(v)-1]
				}
				result.Fl = append(result.Fl, FlagVar{
					F: f,
					V: v,
				})
			}
		}

		if !found {
			return Cmd{}, &UnknownFlagError{S: result.St.Name, F: cur}
		}
	}

	if i < len(token) {
		result.Pa = token[i:]
	}

	return result, nil
}
