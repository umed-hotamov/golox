package interpreter

import (
	"fmt"

	"github.com/umed-hotamov/golox/internal/ast"
)

func (i *Interpreter) execute(statement ast.Stmt) {
  switch statement.(type) {
  case ast.Expression:
    i.executeExpressionStmt(statement.(ast.Expression))
  case ast.Print:
    i.executePrintStmt(statement.(ast.Print))
  case ast.Var:
    i.executeVarStmt(statement.(ast.Var))
  case ast.Block:
    i.executeBlockStmt(statement.(ast.Block), NewEnclosingEnvironment(i.env))
  case ast.If:
    i.executeIfStmt(statement.(ast.If))
  case ast.While:
    i.executeWhileStmt(statement.(ast.While))
  case ast.Function:
    i.executeFunctionStmt(statement.(ast.Function))
  case ast.Return:
    i.executeReturnStmt(statement.(ast.Return))
  }
}

func (i *Interpreter) executeExpressionStmt(statement ast.Expression) {
  i.evaluate(statement.Expression)
}

func (i *Interpreter) executePrintStmt(statement ast.Print) {
  value := i.evaluate(statement.Expression)

  if value == nil {
    fmt.Println("nil")
    return
  }

  fmt.Println(value)
}

func (i *Interpreter) executeVarStmt(statement ast.Var) {
  var value any = nil
  if statement.Initializer != nil {
    value = i.evaluate(statement.Initializer)
  }
  
  i.env.define(statement.Name.Lexeme, value)
}

func (i *Interpreter) executeBlockStmt(statement ast.Block, env *Environment) {
  previous := i.env
  i.env = env

  defer func() {
    i.env = previous
  }()
  
  for _, stmt := range statement.Statements {
    i.execute(stmt)
  } 
}

func (i *Interpreter) executeIfStmt(statement ast.If) {
  if isTruthy(i.evaluate(statement.Condition)) {
    i.execute(statement.ThenBranch)
  } else if statement.ElseBranch != nil {
    i.execute(statement.ElseBranch)
  }
}

func (i *Interpreter) executeWhileStmt(statement ast.While) {
  for isTruthy(i.evaluate(statement.Condition)) {
    i.execute(statement.Body)
  }
}

func (i *Interpreter) executeFunctionStmt(statement ast.Function) {
  function := NewFunction(statement, i.env)
  i.env.define(statement.Name.Lexeme, function)
}

func (i *Interpreter) executeReturnStmt(statement ast.Return) {
  var value any
  if statement.Value != nil {
    value = i.evaluate(statement.Value)
  }

  panic(value)
}
