package interpreter

import (
	"github.com/umed-hotamov/golox/internal/ast"
)

type Callable interface {
  arity() int
  call(interpreter *Interpreter, arguments []any) any
}

type Function struct {
  declaration ast.Function
  closure     *Environment
}

func NewFunction(declaration ast.Function, closure *Environment) *Function {
  return &Function{
    declaration: declaration,
    closure:     closure,
  }
}

func (f *Function) arity() int {
  return len(f.declaration.Params)
}

func (f *Function) call(interpreter *Interpreter, arguments[] any) (value any) {
  env := NewEnclosingEnvironment(f.closure)
  
  for i := 0; i < len(f.declaration.Params); i += 1 {
    env.define(f.declaration.Params[i].Lexeme, arguments[i])
  }
  
  defer func() {
    if r := recover(); r != nil {
       value = r
    }
  }()

  interpreter.executeBlockStmt(f.declaration.Body, env)
  return 
}
