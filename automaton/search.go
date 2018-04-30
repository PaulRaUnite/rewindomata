package automaton

import (
	"unicode/utf8"
	"github.com/uber-go/atomic"
)

//see python realization
func (acc Acceptor) FrontSearch(word string) bool {
	front := acc.initial.union(nil)
	newFront := make(stateSet)

	for _, symb := range word {
		for now := range front {
			delete(front, now)
			if moves, ok := acc.transitions[now]; ok {
				if nexts, ok := moves[symb]; ok {
					for next := range nexts {
						newFront.add(next)
					}
				}
			}
		}
		front, newFront = newFront, front
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

	for ; ; {
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

type optional struct {
	data     payLoad
	found    bool
	presence bool
}

func createWorker(acc Acceptor, in <-chan payLoad, out chan<- optional, termination <-chan struct{}, in_work *atomic.Int64) {
	go func() {
		defer func() {
			v := in_work.Dec()
			if v == 0 {
				close(out)
			}
		}()
		for {
			select {
			case work, ok := <-in:
				{
					if !ok {
						return
					}
					if _, ok := acc.final[work.state]; ok && len(work.rest) == 0 {
						out <- optional{found: true}
						return
					}
					if len(work.rest) != 0 {
						if jumps, ok := acc.transitions[work.state]; ok {
							r, size := utf8.DecodeRuneInString(work.rest)
							if nexts, ok := jumps[r]; ok {
								for state := range nexts {
									out <- optional{data: payLoad{state: state, rest: work.rest[size:]}, presence: true}
								}
							}
						}
					}
					out <- optional{presence: false}
				}
			case <-termination:
				return
			}
		}
	}()
}

//using goroutines
func (acc Acceptor) AtomicParallelSearch(word string, routines int) bool {
	channelCapacity := routines * 8
	if channelCapacity < len(acc.initial) {
		channelCapacity = len(acc.initial)
	}
	queue := make(chan payLoad, channelCapacity*4)
	terminate := make(chan struct{}, 1)
	commonOutput := make(chan optional, channelCapacity)
	in_work := atomic.NewInt64(int64(routines))
	for i := 0; i < routines; i++ {
		createWorker(acc, queue, commonOutput, terminate, in_work)
	}
	for state := range acc.initial {
		queue <- payLoad{state: state, rest: word}
	}
	counter := len(acc.initial)
	found := false

	for value := range commonOutput {
		if value.found {
			found = true
			break
		}
		if value.presence {
			counter++
			queue <- value.data
		} else {
			counter--
			if counter == 0 {
				break
			}
		}
	}
	close(terminate)
	close(queue)
	for v := range commonOutput {
		_ = v
	}
	return found
}

func (acc Acceptor) StochasticSearch(word string) bool {
	queue := newQueue()
	for state := range acc.initial {
		queue.push(payLoad{state: state, rest: word})
	}

	for ; ; {
		work, ok := queue.randomPop()
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
