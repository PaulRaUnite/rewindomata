package automaton

import (
	"unicode/utf8"
	"sync"
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

func merge(cs ...<-chan optional) <-chan optional {
	var wg sync.WaitGroup
	out := make(chan optional, len(cs))

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan optional) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func createWorker(acc Acceptor, in <-chan payLoad, outputSize int) <-chan optional {
	out := make(chan optional, outputSize)
	go func() {
		defer close(out)
		for work := range in {
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
	}()
	return out
}

//using goroutines
func (acc Acceptor) AtomicParallelSearch(word string, routines int) bool {
	channelCapacity := routines * 2
	if channelCapacity < len(acc.initial) {
		channelCapacity = len(acc.initial)
	}
	queue := make(chan payLoad, channelCapacity)
	worker_outputs := make([]<-chan optional, 0, channelCapacity)
	for i := 0; i < routines; i++ {
		worker_outputs = append(worker_outputs, createWorker(acc, queue, channelCapacity))
	}
	for state := range acc.initial {
		queue <- payLoad{state: state, rest: word}
	}
	commonOutput := merge(worker_outputs...)
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
