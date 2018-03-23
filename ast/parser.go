package ast

import (
	st "github.com/golang-collections/collections/stack"
)

func Parse(income string) (AST, error) {
	stack := st.New()

	for i, symbol := range income {
		switch symbol {
		case '|':
			right := stack.Pop()
			if right == nil {
				return AST{}, newError("OR left operand is missed", i)
			}
			left := stack.Pop()
			if left == nil {
				return AST{}, newError("OR right closure operand is missed", i)
			}
			stack.Push(Node{Type: OR, Children: []interface{}{left, right}})
		case '+':
			right := stack.Pop()
			if right == nil {
				return AST{}, newError("AND left operand is missed", i)
			}
			left := stack.Pop()
			if left == nil {
				return AST{}, newError("AND right closure operand is missed", i)
			}
			stack.Push(Node{Type: AND, Children: []interface{}{left, right}})
		case '*':
			operand := stack.Pop()
			if operand == nil {
				return AST{}, newError("closure operand is missed", i)
			}
			stack.Push(Node{Type: CLS, Children: []interface{}{operand}})
		default:
			stack.Push(Leaf{Symbol: symbol})
		}
	}
	return AST{stack.Pop()}, nil
}
