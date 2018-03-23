package automaton

import (
	"github.com/PaulRaUnite/errors"
	"github.com/PaulRaUnite/rewindomata/ast"
)

func NewAcceptor(char rune) Acceptor {
	return Acceptor{
		initial:     stateSet{0: {}},
		final:       stateSet{1: {}},
		transitions: stateTransitions{0: {char: {1: {}}}},
		max:         state(1),
	}
}

func (left Acceptor) And(right Acceptor) Acceptor {
	right = right.shift(left.max + 1)

	newInit := left.initial
	newFinal := right.final
	newTrans := left.transitions.union(right.transitions)

	closures := false
	for state := range left.initial {
		if _, ok := left.final[state]; ok {
			closures = true
			break
		}
	}
	if closures {
		newInit = left.initial.union(right.initial)
	}

	for from, directs := range left.transitions {
		for char, ends := range directs {
			for end := range ends {
				if _, ok := left.final[end]; ok {
					for rightStart := range right.initial {
						newTrans.add(from, char, rightStart)
					}
				}
			}
		}
	}
	return Acceptor{
		initial:     newInit,
		final:       newFinal,
		transitions: newTrans,
		max:         right.max,
	}
}

func (left Acceptor) Or(right Acceptor) Acceptor {
	right = right.shift(left.max + 1)
	newInit := left.initial.union(right.initial)
	newFinal := left.final.union(right.final)
	newTrans := left.transitions.union(right.transitions)

	return Acceptor{
		initial:     newInit,
		final:       newFinal,
		transitions: newTrans,
		max:         right.max,
	}
}

func (left Acceptor) Closure() Acceptor {
	finalStart := left.max + 1

	newInit := stateSet{finalStart: {}}
	newFinal := newInit
	newTrans := make(stateTransitions, len(left.transitions))

	for from, directs := range left.transitions {
		if _, ok := left.initial[from]; ok {
			for char, ends := range directs {
				for end := range ends {
					newTrans.add(finalStart, char, end)
				}
			}
		}
	}
	tailTrans := make(stateTransitions)
	for from, directs := range left.transitions {
		for char, ends := range directs {
			for end := range ends {
				if _, ok := left.final[end]; ok {
					tailTrans.add(from, char, finalStart)
				}
				tailTrans.add(from, char, end)
			}
		}
	}
	newTrans = newTrans.union(tailTrans)

	return Acceptor{
		initial:     newInit,
		final:       newFinal,
		transitions: newTrans,
		max:         finalStart,
	}
}

func ConstructFromAST(ast ast.AST) (Acceptor, error) {
	if ast.Root == nil {
		return Acceptor{initial: stateSet{0: {}}, final: stateSet{0: {}}, transitions: make(stateTransitions), max: 0}, nil
	} else {
		return construct(ast.Root)
	}
}

func construct(root interface{}) (Acceptor, error) {
	if node, ok := root.(ast.Node); ok {
		left, err := construct(node.Children[0])
		if err != nil {
			return left, nil
		}
		if node.Type == ast.OR {
			right, err := construct(node.Children[1])
			if err != nil {
				return right, nil
			}
			return left.Or(right), nil
		} else if node.Type == ast.AND {
			right, err := construct(node.Children[1])
			if err != nil {
				return right, nil
			}
			return left.And(right), nil
		} else if node.Type == ast.CLS {
			return left.Closure(), nil
		} else {
			return Acceptor{}, errors.New("unknown node type")
		}
	} else if leaf, ok := root.(ast.Leaf); ok {
		return NewAcceptor(leaf.Symbol), nil
	} else {
		return Acceptor{}, errors.New("unknown object inside an AST")
	}
}
