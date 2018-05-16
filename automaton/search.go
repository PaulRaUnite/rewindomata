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

type factoryPayload struct {
	state   state
	rest    string
	wordNum int
}
type optional struct {
	data     factoryPayload
	found    bool
	presence bool
}

func createWorker(acc Acceptor, in <-chan factoryPayload, out chan<- optional, termination <-chan struct{}, in_work *atomic.Int64) {
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
					if len(work.rest) == 0 {
						if _, ok := acc.final[work.state]; ok {
							out <- optional{data: work, found: true, presence: false}
						} else {
							out <- optional{found: false, presence: false}
						}
					} else {
						if jumps, ok := acc.transitions[work.state]; ok {
							r, size := utf8.DecodeRuneInString(work.rest)
							if nexts, ok := jumps[r]; ok {
								for state := range nexts {
									out <- optional{
										data: factoryPayload{
											state:   state,
											rest:    work.rest[size:],
											wordNum: work.wordNum,
										},
										found:    false,
										presence: true,
									}
								}
							}
						}
						out <- optional{found: false, presence: false}
					}
				}
			case <-termination:
				return
			}
		}
	}()
}

//using goroutines
func (acc Acceptor) AtomicParallelSearch(words []string, routines int) []bool {
	channelCapacity := len(acc.initial) * len(words) + routines * 8 //??? don't know how to compute channel size

	queue := make(chan factoryPayload, channelCapacity*4)
	terminate := make(chan struct{}, 1)
	commonOutput := make(chan optional, channelCapacity)
	in_work := atomic.NewInt64(int64(routines))
	for i := 0; i < routines; i++ {
		createWorker(acc, queue, commonOutput, terminate, in_work)
	}
	for state := range acc.initial {
		for i, word := range words {
			queue <- factoryPayload{state: state, rest: word, wordNum: i}
		}
	}
	counter := len(acc.initial) * len(words)
	notFound := len(words)

	result := make([]bool, len(words))
	for value := range commonOutput {
		if value.found {
			if result[value.data.wordNum] == false {
				result[value.data.wordNum] = true
				notFound--
				if notFound == 0 {
					break
				}
			}
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
	return result
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
