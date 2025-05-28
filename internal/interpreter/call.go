package interpreter

import (
	"github.com/umed-hotamov/golox/internal/parser"
)

type Callable interface {
  arity() int
  call(interpreter *Interpreter, arguments []any) any
}

type Function struct {
  declaration parser.Function
}

func NewFunction(declaration parser.Function) *Function {
  return &Function{
    declaration: declaration,
  }
}

func (f *Function) arity() int {
  return len(f.declaration.Params)
}

func (f *Function) call(interpreter *Interpreter, arguments[] any) any {
  env := NewEnclosingEnvironment(interpreter.globals)
  
  for i := 0; i < len(f.declaration.Params); i += 1 {
    env.define(f.declaration.Params[i].Lexeme, arguments[i])
  }
  
  interpreter.executeBlockStmt(f.declaration.Body, env)

  return nil
}
