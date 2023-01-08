package core_test

import (
	"core"
	"testing"
)

func TestSearch(t *testing.T) {
	var example = []string{
		"asa",
		"asdf",
		"kajdfoqijfldfzqafsd",
		"fldhgiowfhlkjzlfjq",
		"asfhoqiejfklaj",
		"wattermelon",
		"hotmilk",
		"as it was",
	}
	var ans = []bool{
		true,
		false,
		false,
		false,
		true,
		false,
		false,
		true,
	}

	for i, str := range example {
		if core.Search(str, "as*a") != ans[i] {
			t.Errorf("failed: %s a?s*a %t", str, ans[i])
		}
	}
}
