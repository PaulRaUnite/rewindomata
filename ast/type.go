package ast

import "fmt"

type AST struct {
	Root interface{}
}

type Leaf struct {
	Symbol rune
}

type Node struct {
	Type     NodeType
	Children []interface{} //Leaf or Node
}

type NodeType int

const (
	AND = NodeType(iota)
	OR
	CLS
)

func (l Leaf) String() string {
	return string(l.Symbol)
}

func (n Node) String() string {
	switch n.Type {
	case OR:
		left := n.Children[0].(fmt.Stringer)
		right := n.Children[1].(fmt.Stringer)
		return fmt.Sprintf("%s%s%s", left.String(), right.String(), "|")
	case AND:
		left := n.Children[0].(fmt.Stringer)
		right := n.Children[1].(fmt.Stringer)
		return fmt.Sprintf("%s%s%s", left.String(), right.String(), "+")
	case CLS:
		operand := n.Children[0].(fmt.Stringer)
		return fmt.Sprintf("%s%s", operand, "*")
	default:
		panic("there is no other values in the enum")
	}
}

func (ast AST) String() string {
	if ast.Root == nil {
		return ""
	} else {
		if node, ok := ast.Root.(fmt.Stringer); ok {
			return node.String()
		} else {
			panic("ast.Root contains undetected object")
		}
	}
}
