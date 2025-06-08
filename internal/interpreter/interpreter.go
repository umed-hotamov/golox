package interpreter

import (
	"fmt"

	"github.com/umed-hotamov/golox/internal/lexer"
	"github.com/umed-hotamov/golox/internal/ast"
)

type Interpreter struct {
  env     *Environment
  globals *Environment
}

func NewInterpreter() *Interpreter {
  globals := NewEnvironment()

  globals.define("clock", new(Clock))

  return &Interpreter{
    env:     globals,
    globals: globals,
  }
}

func (i *Interpreter) Interpret(statements []ast.Stmt) {
  defer errorRecovery()

  for _, stmt := range statements {
    i.execute(stmt)
  }
}

func isTruthy(value any) bool {
  if value == nil {
    return false
  }

  if isNumber(value) {
    return value.(float64) != 0
  }
  if isString(value) {
    return value.(string) != ""
  }
  
  if isBool(value) {
    return value.(bool)
  }

  return true
}

func isString(value any) bool {
  _, ok := value.(string)
  return ok
}

func isNumber(value any) bool {
  _, ok := value.(float64)
  return ok
}

func isBool(value any) bool {
  _, ok := value.(bool)
  return ok
}

func isEqual(left any, right any) bool {
  if left == nil && right == nil {
    return true
  }
  if left == nil {
    return false
  }

  if isString(left) && isString(right) {
    return left.(string) == right.(string)
  }
  if isNumber(left) && isNumber(right) {
    return left.(float64) == right.(float64)
  }

  return false
}

func runtimeError(token lexer.Token, message string) {
  panic(fmt.Sprintf("[line: %d , at %s] Error: %s\n", token.Line, token.Lexeme, message))
}

func errorRecovery() {
  if err := recover(); err != nil {
    fmt.Print(err)
  }
}
