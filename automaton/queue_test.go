package automaton

import (
	"testing"
)

func TestQueue(t *testing.T) {
	e1 := payLoad{1, "a"}
	e2 := payLoad{2, "b"}
	q := newQueue()
	q.push(e1)
	if out, ok := q.pop(); ok {
		if out != e1 {
			t.Fail()
		}
	} else {
		t.Fail()
	}
	q.push(e2)
	q.push(e1)
	if out, ok := q.pop(); ok {
		if out != e2 || q.len() != 1 {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}

func TestRandomQueue(t *testing.T) {
	e1 := payLoad{1, "a"}
	e2 := payLoad{2, "b"}
	q := newQueue()
	q.push(e1)
	q.push(e2)
	if out, ok := q.randomPop(); ok {
		if q.len() != 1 {
			t.Fail()
		}
		if out != e1 || out != e2 {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}
