package resolver

import "github.com/umed-hotamov/golox/internal/ast"

func (r *Resolver) resolveStatement(statement ast.Stmt) {
	switch statement.(type) {
	case ast.Block:
		r.resolveBlock(statement.(ast.Block))
	case ast.Var:
		r.resolveVar(statement.(ast.Var))
	case ast.Function:
		r.resolveFunction(statement.(ast.Function))
	case ast.Expression:
		r.resolveExpressionStatement(statement.(ast.Expression))
	case ast.Print:
		r.resolvePrint(statement.(ast.Print))
	case ast.Return:
		r.resolveReturn(statement.(ast.Return))
	case ast.While:
		r.resolveWhile(statement.(ast.While))
	case ast.Class:
		r.resolveClass(statement.(ast.Class))
	}
}

func (r *Resolver) resolveBlock(statement ast.Block) {
	r.beginScope()
	r.Resolve(statement.Statements)
	r.endScope()
}

func (r *Resolver) resolveVar(statement ast.Var) {
	r.declare(statement.Name)
	if statement.Initializer != nil {
		r.resolveExpression(statement.Initializer)
	}

	r.define(statement.Name)
}

func (r *Resolver) resolveFunction(statement ast.Function) {
	r.declare(statement.Name)
	r.define(statement.Name)

	enclosingFunction := r.currentFunction
	r.currentFunction = FUNCTION

	r.beginScope()
	for _, param := range statement.Params {
		r.declare(param)
		r.define(param)
	}
	r.Resolve(statement.Body.Statements)
	r.endScope()

	r.currentFunction = enclosingFunction
}

func (r *Resolver) resolveExpressionStatement(statement ast.Expression) {
	r.resolveExpression(statement.Expression)
}

func (r *Resolver) resolveIf(statement ast.If) {
	r.resolveExpression(statement.Condition)
	r.resolveStatement(statement.ThenBranch)

	if statement.ElseBranch != nil {
		r.resolveExpression(statement.ElseBranch)
	}
}

func (r *Resolver) resolvePrint(statement ast.Print) {
	r.resolveExpression(statement.Expression)
}

func (r *Resolver) resolveReturn(statement ast.Return) {
	if r.currentFunction == NONE {
		r.error(statement.Keyword, "Can't return from top-level code")
	}

	if statement.Value != nil {
		r.resolveExpression(statement.Value)
	}
}

func (r *Resolver) resolveWhile(statement ast.While) {
	r.resolveExpression(statement.Condition)
	r.resolveStatement(statement.Body)
}

func (r *Resolver) resolveClass(statement ast.Class) {
	r.declare(statement.Name)
	r.define(statement.Name)
}
