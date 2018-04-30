package automaton

type queue struct {
	start int
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
	defer func() {q.start += 1}()
	return q.buffer[q.start], true
}

func (q queue) len() int {
	return len(q.buffer) - q.start
}