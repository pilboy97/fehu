package core

import (
	"regexp"
	"time"

	"github.com/pkg/errors"
)

var TimeFmt = `2006-01-02;15:04:05`
var ErrWrongName = errors.New("wrong name")
var ErrCannotFind = errors.New("failed to find")
var ErrAlreadyExists = errors.New("already exists")

func SureID(ret int64) int64 {
	switch ret {
	case -1:
		panic(ErrCannotFind)
	case -2:
		panic(ErrAlreadyExists)
	default:
		break
	}

	return ret
}

func SureName(name string) string {
	ok, err := regexp.MatchString(`^~?((\p{L}|_)(\p{L}|\d|_)*:)*((\p{L}|_)(\p{L}|\d|_)*)$`, name)
	if err != nil {
		panic(err)
	}

	if !ok {
		panic(ErrWrongName)
	}

	if id := GetAccByName(name); id != -1 {
		panic(ErrAlreadyExists)
	}

	return name
}
func Search(str, ptn string) bool {
	var D = make([][]bool, len(ptn))
	for i := range D {
		D[i] = make([]bool, len(str))
	}

	//fmt.Println(str, ptn)
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

	/*
		for i := 0; i < len(P); i++ {
			for j := 0; j < len(S); j++ {
				if D[i][j] {
					fmt.Print("+")
				} else {
					fmt.Print("-")
				}
			}
			fmt.Println()
		}
	*/

	for i := 0; i < len(S); i++ {
		if D[len(P)-1][i] {
			return true
		}
	}

	return false
}
func ParseTime(str string) time.Time {
	ret, err := time.ParseInLocation(
		TimeFmt,
		str,
		time.Local)

	if err != nil {
		panic(err)
	}

	return ret
}

func ParsePeriod(st, ed string) Period {
	var A, B *time.Time

	if len(st) != 0 {
		A = &time.Time{}
		*A = ParseTime(st)
	}
	if len(ed) != 0 {
		B = &time.Time{}
		*B = ParseTime(ed)
	}

	return Period{St: A, Ed: B}
}
