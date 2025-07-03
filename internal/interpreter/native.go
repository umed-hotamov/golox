package interpreter

import "time"

type Clock struct {
}

func (c Clock) arity() int {
	return 0
}

func (c Clock) call(interpreter *Interpreter, arguments []any) any {
	return float64(time.Now().Second())
}
