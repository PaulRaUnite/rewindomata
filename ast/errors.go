package ast

import (
	"fmt"
)

type Error struct {
	msg    string
	offset int
}

func newError(msg string, offset int) Error {
	return Error{msg, offset}
}

func (p Error) Error() string {
	return fmt.Sprintf("parse error at position %d: %s", p.offset, p.msg)
}
