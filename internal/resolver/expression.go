package resolver

import "github.com/umed-hotamov/golox/internal/ast"

func (r *Resolver) resolveExpression(expression ast.Expr) {
	switch expression.(type) {
	case ast.Variable:
    r.resolveVariable(expression.(ast.Variable))
	}
}

func (r *Resolver) resolveVariable(expression ast.Variable) {
	if !r.scopes.IsEmpty() {
		scope := r.scopes.Peek().(map[string]bool)
		if v, ok := scope[expression.Name.Lexeme]; ok {
			if !v {
				r.error(expression.Name, "Can't read local variable in its own initializer")
			}
		}
	}

	r.resolveLocal(expression, expression.Name)
}
