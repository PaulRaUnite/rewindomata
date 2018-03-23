package automaton

import (
	"testing"
	"github.com/PaulRaUnite/rewindomata/ast"
)

func TestAcceptorBuilder_And(t *testing.T) {
	tree, err := ast.Parse("ab+")
	if err != nil {
		t.Fatal(err)
	}
	acc,err := ConstructFromAST(tree)
	if err != nil {
		t.Fatal(err)
	}
	if acc.FrontSearch("ab") != true {
		t.Fatal()
	}
}

func TestAcceptorBuilder_Or(t *testing.T) {
	tree, err := ast.Parse("ab|")
	if err != nil {
		t.Fatal(err)
	}
	acc,err := ConstructFromAST(tree)
	if err != nil {
		t.Fatal(err)
	}
	if acc.FrontSearch("a") != true {
		t.Fatal()
	}
	if acc.FrontSearch("b") != true {
		t.Fatal()
	}
}

func TestAcceptorBuilder_Closure(t *testing.T) {
	tree, err := ast.Parse("ab+*")
	if err != nil {
		t.Fatal(err)
	}
	acc,err := ConstructFromAST(tree)
	if err != nil {
		t.Fatal(err)
	}
	if acc.FrontSearch("") != true {
		t.Fatal()
	}
	if acc.FrontSearch("ab") != true {
		t.Fatal()
	}
	if acc.FrontSearch("abab") != true {
		t.Fatal()
	}
	if acc.FrontSearch("aba") != false {
		t.Fatal()
	}
	if acc.FrontSearch("b") != false {
		t.Fatal()
	}
}
