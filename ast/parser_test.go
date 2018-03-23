package ast

import (
	"testing"
)

func TestParse(t *testing.T) {
	in := "ab+c|*"
	ast, err := Parse(in)
	if err != nil {
		t.Fatal(err)
	}

	if astStr := ast.String(); astStr != in {
		t.Fatal("input doesn't equal to ast stringification:", astStr, "!=", in)
	}
}
