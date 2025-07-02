package resolver

import (
	"fmt"

	"github.com/umed-hotamov/golox/internal/ast"
	"github.com/umed-hotamov/golox/internal/interpreter"
	"github.com/umed-hotamov/golox/internal/lexer"
)

type Resolver struct {
	interpreter *interpreter.Interpreter
	scopes      *Stack
	hasError    bool
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

	scope := r.scopes.Peek().(map[string]bool)
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name lexer.Token) {
	if r.scopes.IsEmpty() {
		return
	}

	scope := r.scopes.Peek().(map[string]bool)
	scope[name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expression ast.Expr, name lexer.Token) {
	for i := r.scopes.Size() - 1; i >= 0; i-- {
		scope := r.scopes.Get(i).(map[string]bool)
		if _, ok := scope[name.Lexeme]; ok {
			r.interpreter.Resolve(expression, r.scopes.Size()-1-i)
			return
		}
	}
}

func (r *Resolver) error(token lexer.Token, message string) {
	fmt.Printf("[line: %d] Error: %s\n", token.Line, message)
	r.hasError = true
}
