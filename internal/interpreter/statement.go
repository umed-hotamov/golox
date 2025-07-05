package interpreter

import (
	"fmt"

	"github.com/umed-hotamov/golox/internal/ast"
)

func (i *Interpreter) execute(statement ast.Stmt) {
	switch statement.(type) {
	case ast.Expression:
		i.executeExpression(statement.(ast.Expression))
	case ast.Print:
		i.executePrint(statement.(ast.Print))
	case ast.Var:
		i.executeVar(statement.(ast.Var))
	case ast.Block:
		i.executeBlock(statement.(ast.Block), NewEnclosingEnvironment(i.env))
	case ast.If:
		i.executeIf(statement.(ast.If))
	case ast.While:
		i.executeWhile(statement.(ast.While))
	case ast.Function:
		i.executeFunction(statement.(ast.Function))
	case ast.Return:
		i.executeReturn(statement.(ast.Return))
  case ast.Class:
    i.executeClass(statement.(ast.Class))
  }
}

func (i *Interpreter) executeExpression(statement ast.Expression) {
	i.evaluate(statement.Expression)
}

func (i *Interpreter) executePrint(statement ast.Print) {
	value := i.evaluate(statement.Expression)

	if value == nil {
		fmt.Println("nil")
		return
	}

	fmt.Println(value)
}

func (i *Interpreter) executeVar(statement ast.Var) {
	var value any = nil
	if statement.Initializer != nil {
		value = i.evaluate(statement.Initializer)
	}

	i.env.define(statement.Name.Lexeme, value)
}

func (i *Interpreter) executeBlock(statement ast.Block, env *Environment) {
	previous := i.env
	i.env = env

	defer func() {
		i.env = previous
	}()

	for _, stmt := range statement.Statements {
		i.execute(stmt)
	}
}

func (i *Interpreter) executeIf(statement ast.If) {
	if isTruthy(i.evaluate(statement.Condition)) {
		i.execute(statement.ThenBranch)
	} else if statement.ElseBranch != nil {
		i.execute(statement.ElseBranch)
	}
}

func (i *Interpreter) executeWhile(statement ast.While) {
	for isTruthy(i.evaluate(statement.Condition)) {
		i.execute(statement.Body)
	}
}

func (i *Interpreter) executeFunction(statement ast.Function) {
	function := NewFunction(statement, i.env)
	i.env.define(statement.Name.Lexeme, function)
}

func (i *Interpreter) executeReturn(statement ast.Return) {
	var value any
	if statement.Value != nil {
		value = i.evaluate(statement.Value)
	}

	panic(value)
}

func (i *Interpreter) executeClass(statement ast.Class) {
  i.env.define(statement.Name.Lexeme, nil)
  class := NewLoxClass(statement.Name.Lexeme)
  i.env.assign(statement.Name, class)
}
