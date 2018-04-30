package automaton

import (
	"testing"
	"fmt"
	"github.com/PaulRaUnite/rewindomata/ast"
)

const ROUTINES = 4

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
				if acc.AtomicParallelSearch(c, ROUTINES) != result {
					fmt.Println(acc)
					t.Error(regexp, result, c)
				}
			}
		}
	}
}

func benchSearches(b *testing.B, result bool, searchFunc func(acc Acceptor, word string) bool) {
	b.StopTimer()
	for regexp, resultCases := range examples {
		tree, err := ast.Parse(regexp)
		if err != nil {
			b.Fatal(err)
		}
		acc, err := ConstructFromAST(tree)
		if err != nil {
			b.Fatal(err)
		}
		for _, c := range resultCases[result] {
			b.StartTimer()
			for n := 0; n < b.N; n++ {
				if searchFunc(acc, c) != result {
					fmt.Println(acc)
					b.Error(regexp, result, c)
				}
			}
			b.StopTimer()
		}
	}
	b.StartTimer()
}
func BenchmarkAcceptor_FrontSearchPositive(b *testing.B) {
	benchSearches(b, true, func(acc Acceptor, word string) bool {
		return acc.FrontSearch(word)
	})
}

func BenchmarkAcceptor_AtomicSearchPositive(b *testing.B) {
	benchSearches(b, true, func(acc Acceptor, word string) bool {
		return acc.AtomicSearch(word)
	})
}

func BenchmarkAcceptor_AtomicParallelSearchPositive(b *testing.B) {
	benchSearches(b, true, func(acc Acceptor, word string) bool {
		return acc.AtomicParallelSearch(word, ROUTINES)
	})
}

func BenchmarkAcceptor_StochasticSearchPositive(b *testing.B) {
	benchSearches(b, true, func(acc Acceptor, word string) bool {
		return acc.StochasticSearch(word)
	})
}

func BenchmarkAcceptor_FrontSearchNegative(b *testing.B) {
	benchSearches(b, false, func(acc Acceptor, word string) bool {
		return acc.FrontSearch(word)
	})
}

func BenchmarkAcceptor_AtomicSearchNegative(b *testing.B) {
	benchSearches(b, false, func(acc Acceptor, word string) bool {
		return acc.AtomicSearch(word)
	})
}

func BenchmarkAcceptor_AtomicParallelSearchNegative(b *testing.B) {
	benchSearches(b, false, func(acc Acceptor, word string) bool {
		return acc.AtomicParallelSearch(word, ROUTINES)
	})
}

func BenchmarkAcceptor_StochasticSearchNegative(b *testing.B) {
	benchSearches(b, false, func(acc Acceptor, word string) bool {
		return acc.StochasticSearch(word)
	})
}
