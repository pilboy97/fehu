package core

import (
	"regexp"
	"time"

	"github.com/pkg/errors"
)

var TimeFmt = `2006-01-02;15:04:05`

func ErrWrongName(name string) error {
	return errors.Errorf("wrong name: %s", name)
}
func ErrCannotFind(name string) error {
	return errors.Errorf("cannot find: %s", name)
}
func ErrAlreadyExists(name string) error {
	return errors.Errorf("already exists: %s", name)
}

func SureName(name string) (string, error) {
	ok, err := regexp.MatchString(`^~?((\p{L}|_)(\p{L}|\d|_)*:)*((\p{L}|_)(\p{L}|\d|_)*)$`, name)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", ErrWrongName(name)
	}

	// Check if it's a reserved variable name
	if _, ok := Vars[name]; ok {
		return "", ErrAlreadyExists(name)
	}

	// Check if an account with this name already exists
	if _, err := GetAccByName(name); err == nil { // No error means it exists
		return "", ErrAlreadyExists(name)
	}

	// If all checks pass, the name is sure
	return name, nil
}
func Search(str, ptn string) bool {
	var D = make([][]bool, len(ptn))
	for i := range D {
		D[i] = make([]bool, len(str))
	}

	S, P := []rune(str), []rune(ptn)

	if S[0] == P[0] || P[0] == '*' || P[0] == '?' {
		D[0][0] = true
	}

	for j := 1; j < len(S); j++ {
		if P[0] == '*' {
			D[0][j] = D[0][j-1]
		}
	}
	for i := 1; i < len(P); i++ {
		for j := 1; j < len(S); j++ {
			if P[i] == S[j] {
				D[i][j] = D[i-1][j-1]
			} else if P[i] == '?' {
				D[i][j] = D[i-1][j-1]
			} else if P[i] == '*' {
				D[i][j] = (D[i][j-1] || D[i-1][j] || D[i-1][j-1])
			}
		}
	}

	for i := 0; i < len(S); i++ {
		if D[len(P)-1][i] {
			return true
		}
	}

	return false
}
func ParseTime(str string) (int64, error) {
	ret, err := time.ParseInLocation(
		TimeFmt,
		str,
		time.Local)

	if err != nil {
		return 0, err
	}

	return ret.UTC().Unix(), nil
}

func ParsePeriod(st, ed string) Period {
	var A, B *int64

	if len(st) != 0 {
		ts, err := ParseTime(st)
		if err != nil {
			// Handle error, perhaps return error from ParsePeriod
		}
		A = &ts
	}
	if len(ed) != 0 {
		ts, err := ParseTime(ed)
		if err != nil {
			// Handle error, perhaps return error from ParsePeriod
		}
		B = &ts
	}

	return Period{St: A, Ed: B}
}
