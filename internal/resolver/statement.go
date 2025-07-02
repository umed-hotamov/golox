package resolver

import "github.com/umed-hotamov/golox/internal/ast"

func (r *Resolver) resolveStatement(statement ast.Stmt) {
	switch statement.(type) {
	case ast.Block:
		r.resolveBlock(statement.(ast.Block))
	}
}

func (r *Resolver) resolveBlock(statement ast.Block) {
	r.beginScope()
	r.resolve(statement.Statements)
	r.endScope()
}

func (r *Resolver) resolveVar(statement ast.Var) {
	r.declare(statement.Name)
	if statement.Initializer != nil {
		r.resolveExpression(statement.Initializer)
	}

	r.define(statement.Name)
}
