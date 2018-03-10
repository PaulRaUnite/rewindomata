package automata

// state type used for type safe
// manipulations
type state uint64

type stateSet map[state]struct{}

func (ss stateSet) add(value state) {
	ss[value] = struct{}{}
}

func (ss stateSet) union(set stateSet) stateSet {
	newSet := make(stateSet, len(ss)+len(set))
	for x := range ss {
		newSet[x] = struct{}{}
	}
	for x := range set {
		newSet[x] = struct{}{}
	}
	return newSet
}


func (st stateTransitions) add(from state, char rune, to state) {
	if dir, ok := st[from]; ok {
		if s, ok := dir[char]; ok {
			s.add(to)
		} else {
			dir[char] = stateSet{to: {}}
		}
	} else {
		st[from] = map[rune]stateSet{char: {to: {}}}
	}
}

type stateTransitions map[state]map[rune]stateSet

func (left stateTransitions) union(right stateTransitions) stateTransitions {
	newTrans := make(stateTransitions, len(left)+len(right))

	for from, direction := range left {
		for char, ends := range direction {
			for to := range ends {
				newTrans.add(from, char, to)
			}
		}
	}
	for from, direction := range right {
		for char, ends := range direction {
			for to := range ends {
				newTrans.add(from, char, to)
			}
		}
	}
	return newTrans
}

type Acceptor struct {
	initial     stateSet
	final       stateSet
	transitions stateTransitions
	max         state
}

func (acc Acceptor) shift(by state) Acceptor {
	newInit := make(stateSet, len(acc.initial))
	newFinal := make(stateSet, len(acc.final))
	newTrans := make(stateTransitions, len(acc.transitions))

	for x := range acc.initial {
		newInit[x+by] = struct{}{}
	}

	for x := range acc.final {
		newFinal[x+by] = struct{}{}
	}

	for from, direction := range acc.transitions {
		for char, ends := range direction {
			for to := range ends {
				newTrans.add(from+by, char, to+by)
			}
		}
	}
	return Acceptor{
		initial:     newInit,
		final:       newFinal,
		transitions: newTrans,
		max:         acc.max + by,
	}
}
