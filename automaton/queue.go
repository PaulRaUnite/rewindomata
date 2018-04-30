package automaton

import "math/rand"

type queue struct {
	start  int
	buffer []payLoad
}

func newQueue() *queue {
	return &queue{start: 0, buffer: make([]payLoad, 0, 16)}
}

func (q *queue) push(value payLoad) {
	q.buffer = append(q.buffer, value)
}
func (q *queue) pop() (payLoad, bool) {
	if q.len() == 0 {
		q.start = 0
		q.buffer = q.buffer[:0]
		return payLoad{}, false
	}
	value := q.buffer[q.start]
	q.start += 1
	return value, true
}

func (q *queue) randomPop() (payLoad, bool) {
	length := q.len()
	if q.len() == 0 {
		return payLoad{}, false
	}
	i := rand.Intn(length)
	value := q.buffer[i]
	q.buffer[i] = q.buffer[length-1]
	q.buffer = q.buffer[:length-1]
	return value, true
}

func (q queue) len() int {
	return len(q.buffer) - q.start
}
