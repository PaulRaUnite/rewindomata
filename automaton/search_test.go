package automaton

import (
	"fmt"
	"testing"

	"github.com/PaulRaUnite/rewindomata/ast"
)

func TestAcceptor_Searches(t *testing.T) {
	examples := map[string]map[bool][]string{
		"ab+":    {true: {"ab"}, false: {"a", "b"}},
		"ab|":    {true: {"a", "b"}, false: {"ab"}},
		"ab+*":   {true: {"", "ab", "ab"}, false: {"ba", "bb"}},
		"abc++":  {true: {"abc"}, false: {"ba", "bb", "ac"}},
		"ab+c*|": {true: {"", "cc", "ab"}, false: {"ac", "g"}},
	}

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
					t.Fatal(regexp, result, c)
				}
				if acc.FrontSearch(c) != result {
					fmt.Println(acc)
					t.Fatal(regexp, result, c)
				}
				if acc.AtomicParallelSearch(c, 5) != result {
					fmt.Println(acc)
					t.Fatal(regexp, result, c)
				}
			}
		}
	}
}
