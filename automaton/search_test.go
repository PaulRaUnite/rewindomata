package automaton

import (
	"testing"
	"fmt"
	"github.com/PaulRaUnite/rewindomata/ast"
)

var examples = map[string]map[bool][]string{
	"ab+":    {true: {"ab"}, false: {"a", "b", ""}},
	"ab|":    {true: {"a", "b"}, false: {"ab", "bc", "ac"}},
	"ab+*":   {true: {"", "ab", "abab", "ababab"}, false: {"ba", "bb", "aa", "cc"}},
	"abc++":  {true: {"abc"}, false: {"ba", "bb", "ac"}},
	"ab+c*|": {true: {"", "cc", "ab"}, false: {"ac", "g"}},
	"ab+c|":  {true: {"ab", "c"}, false: {"abc"}},
	"a*":     {true: {"a", "aa", "aaa"}, false: {"b", "c", "abb"}},
	"ab|*":   {true: {"aabb", "abba", "abaa"}, false: {"aaac", "bbbbbbbbbbad"}},
}

func TestAcceptor_Searches(t *testing.T) {
	for regexp, resultCases := range examples {
		tree, err := ast.Parse(regexp)
		if err != nil {
			t.Fatal(err)
		}
		acc, err := ConstructFromAST(tree)
		if err != nil {
			t.Fatal(err)
		}
		for result, cases := range resultCases {
			for _, c := range cases {
				if acc.AtomicSearch(c) != result {
					fmt.Println(acc)
					t.Error(regexp, result, c)
				}
				if acc.FrontSearch(c) != result {
					fmt.Println(acc)
					t.Error(regexp, result, c)
				}
				if acc.StochasticSearch(c) != result {
					fmt.Println(acc)
					t.Error(regexp, result, c)
				}
			}
		}
	}
}
