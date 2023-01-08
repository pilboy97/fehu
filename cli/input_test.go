package cli_test

import (
	"cli"
	"reflect"
	"testing"
)

func TestTokenizeCommand(t *testing.T) {
	var cases = []struct {
		Input  string
		Output [][]rune
	}{
		{
			`aaaa bbbb`,
			[][]rune{[]rune("aaaa"), []rune("bbbb")},
		},
		{
			`aaa bbbb -c=dd`,
			[][]rune{[]rune("aaa"), []rune("bbbb"), []rune("-c=dd")},
		},
		{
			`"a" "bb"  -c`,
			[][]rune{[]rune("\"a\""), []rune("\"bb\""), []rune("-c")},
		},
	}

	for i, C := range cases {
		res := cli.TokenizeCommand(C.Input)

		if !reflect.DeepEqual(res, C.Output) {
			t.Errorf("\ncase #%2d:\n\t%v,\n\t%v", i, res, C.Output)
		}
	}
}
