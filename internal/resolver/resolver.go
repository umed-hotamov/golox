package resolver

import (
	"github.com/umed-hotamov/golox/internal/ast"
	"github.com/umed-hotamov/golox/internal/interpreter"
	"github.com/umed-hotamov/golox/internal/lexer"
)

type Resolver struct {
  interpreter *interpreter.Interpreter
  scopes      *Stack
}

func NewResolver(interpreter *interpreter.Interpreter) *Resolver {
  return &Resolver{
    interpreter: interpreter,
    scopes:      NewStack(),
  }
}

func (r *Resolver) resolve(statements []ast.Stmt) {
  for _, statement := range statements {
    r.resolveStatement(statement)
  }
}

func (r *Resolver) beginScope() {
  r.scopes.Push(make(map[string]bool))
}

func (r *Resolver) endScope() {
  r.scopes.Pop()
}

func (r *Resolver) declare(name lexer.Token) {
  if r.scopes.IsEmpty() {
    return
  }

  scope := r.scopes.Peek()
  scope.(map[string]bool)[name.Lexeme] = false
}

func (r *Resolver) define(name lexer.Token) {
  if r.scopes.IsEmpty() {
    return
  }

  scope := r.scopes.Peek()
  scope.(map[string]bool)[name.Lexeme] = true
}
