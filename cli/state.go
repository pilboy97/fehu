package cli

import (
	"fmt"
	"regexp"
	"strings"
)

type State struct {
	next []*State

	Action bool
	Name   string
	Manual string
	Pat    string
	Param  int
	Flags  []*Flag
}

func (s *State) SetNext(n *State) {
	s.next = append(s.next, n)
}
func (s *State) HelpString() string {
	result := make([]string, 0, len(s.next)+1)
	result = append(result, s.Name)
	result = append(result, s.Manual)

	for _, child := range s.next {
		result = append(result, fmt.Sprintf("\t%s", child.Name))
		result = append(result, fmt.Sprintf("\t%s", child.Manual))
	}

	return strings.Join(result, "\n")
}
func (s *State) Match(str string) (bool, error) {
	return regexp.MatchString(s.Pat, str)
}

type Flag struct {
	Name    string
	Manual  string
	NamePat string
	ValPat  string
}

func (f *Flag) Match(str string) (bool, error) {
	token := ParseFlag(str)
	if token == nil {
		return false, nil
	}

	ok, err := regexp.MatchString(f.NamePat, token[0])
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	switch len(token) {
	case 1:
		return true, err
	case 2:
		ok, err = regexp.MatchString(f.ValPat, token[1])
		return ok, err
	}

	return true, nil
}
