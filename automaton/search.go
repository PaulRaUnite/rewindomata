package automaton

import (
	"unicode/utf8"
)

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

type payLoad struct {
	state state
	rest  string
}
//see automaton lectures
func (acc Acceptor) AtomicSearch(word string) bool {
	queue := newQueue()
	for state := range acc.initial {
		queue.push(payLoad{state: state, rest: word})
	}

	for ;; {
		work, ok := queue.pop()
		if ok != true {
			break
		}
		if _, ok := acc.final[work.state]; ok && len(work.rest) == 0 {
			return true
		}
		if len(work.rest) == 0 {
			continue
		}
		if jumps, ok := acc.transitions[work.state]; ok {
			r, size := utf8.DecodeRuneInString(work.rest)
			if nexts, ok := jumps[r]; ok {
				for state := range nexts {
					queue.push(payLoad{state: state, rest: work.rest[size:]})
				}
			}
		}
	}
	return false
}

//using goroutines
func (acc Acceptor) AtomicParallelSearch(word string, routines int) bool {
	queue := make(chan payLoad, len(acc.initial)*2)
	termination := make(chan struct{}, 1)
	output := make(chan struct{}, 1)
	scores := make(chan int, routines)

	for state := range acc.initial {
		queue <- payLoad{state: state, rest: word}
	}
	worker := func(queue chan payLoad, scores chan<- int, output chan<- struct{}, terminate <-chan struct{}) {
		for {
			select {
			case work := <-queue:
				if len(work.rest) == 0 {
					if _, ok := acc.final[work.state]; ok {
						output <- struct{}{}
						return
					} else {
						scores <- -1
						continue
					}
				}
				if jumps, ok := acc.transitions[work.state]; ok {
					r, size := utf8.DecodeRuneInString(work.rest)
					if nexts, ok := jumps[r]; ok {
						for state := range nexts {
							queue <- payLoad{state: state, rest: work.rest[size:]}
						}

						scores <- len(nexts)
					}
				}
				scores <- -1
			case <-terminate:
				return
			}
		}
	}
	for i := 0; i < routines; i++ {
		go worker(queue, scores, output, termination)
	}

	dataSize := len(acc.initial)
	for {
		select {
		case <-output:
			close(termination)
			return true
		case increase := <-scores:
			dataSize += increase
			if dataSize == 0 {
				close(termination)
				return false
			}
		}
	}
	return false
}
