package resolver

import "github.com/umed-hotamov/golox/internal/ast"

func (r *Resolver) resolveExpression(expression ast.Expr) {
	switch expression.(type) {
	case ast.Variable:
		r.resolveVariable(expression.(ast.Variable))
	case ast.Assign:
		r.resolveAssign(expression.(ast.Assign))
	case ast.Binary:
		r.resolveBinary(expression.(ast.Binary))
	case ast.Call:
		r.resolveCall(expression.(ast.Call))
	case ast.Grouping:
		r.resolveGrouping(expression.(ast.Grouping))
	case ast.Literal:
		r.resolveLiteral(expression.(ast.Literal))
	case ast.Unary:
		r.resolveUnary(expression.(ast.Unary))
	case ast.Logical:
		r.resolveLogical(expression.(ast.Logical))
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

func (r *Resolver) resolveAssign(expression ast.Assign) {
	r.resolveExpression(expression.Value)
	r.resolveLocal(expression, expression.Name)
}

func (r *Resolver) resolveBinary(expression ast.Binary) {
	r.resolveExpression(expression.Left)
	r.resolveExpression(expression.Right)
}

func (r *Resolver) resolveCall(expression ast.Call) {
	r.resolveExpression(expression.Callee)

	for _, arg := range expression.Arguments {
		r.resolveExpression(arg)
	}
}

func (r *Resolver) resolveGrouping(expression ast.Grouping) {
	r.resolveExpression(expression.Expr)
}

func (r *Resolver) resolveLiteral(expression ast.Literal) {
	_ = expression
}

func (r *Resolver) resolveUnary(expression ast.Unary) {
	r.resolveExpression(expression.Right)
}

func (r *Resolver) resolveLogical(expression ast.Logical) {
	r.resolveExpression(expression.Left)
	r.resolveExpression(expression.Right)
}
