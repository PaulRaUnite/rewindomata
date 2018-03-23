package automaton

//see python realization
func (acc Acceptor) FrontSearch(word string) bool {
	front := acc.initial

	for _, symb := range word {
		newFront := make(stateSet)
		for now := range front {
			if moves, ok := acc.transitions[now]; ok {
				if nexts, ok := moves[symb]; ok {
					for next := range nexts {
						newFront.add(next)
					}
				}
			}
		}
		front = newFront
	}
	for state := range front {
		if _, ok := acc.final[state]; ok {
			return true
		}
	}
	return false
}