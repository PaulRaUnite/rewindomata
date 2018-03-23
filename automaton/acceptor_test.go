package automaton

import (
	"reflect"
	"testing"
)

func TestSetAdd(t *testing.T) {
	examples := []struct {
		left    stateSet
		add     []state
		outcome stateSet
	}{
		{left: stateSet{0: {}}, add: []state{1, 2}, outcome: stateSet{1: {}, 0: {}, 2: {}}},
	}
	for i, example := range examples {
		for _, state := range example.add {
			example.left.add(state)
		}
		if reflect.DeepEqual(example.left, example.outcome) == false {
			panic(i)
		}
	}
}

func TestSetUnion(t *testing.T) {
	examples := []struct {
		left    stateSet
		right   stateSet
		outcome stateSet
	}{
		{left: stateSet{0: {}}, right: stateSet{1: {}}, outcome: stateSet{1: {}, 0: {}}},
	}
	for i, example := range examples {
		if reflect.DeepEqual(example.left.union(example.right), example.outcome) == false {
			panic(i)
		}
	}
}
