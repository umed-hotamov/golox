package interpreter

import (
	"fmt"

	"github.com/umed-hotamov/golox/internal/parser"
)

func (i *Interpreter) execute(statement parser.Stmt) {
  switch statement.(type) {
  case parser.Expression:
    i.executeExpressionStmt(statement.(parser.Expression))
  case parser.Print:
    i.executePrintStmt(statement.(parser.Print))
  case parser.Var:
    i.executeVarStmt(statement.(parser.Var))
  case parser.Block:
    i.executeBlockStmt(statement.(parser.Block), NewEnclosingEnvironment(i.env))
  case parser.If:
    i.executeIfStmt(statement.(parser.If))
  case parser.While:
    i.executeWhileStmt(statement.(parser.While))
  }
}

func (i *Interpreter) executeExpressionStmt(statement parser.Expression) {
  i.evaluate(statement.Expression)
}

func (i *Interpreter) executePrintStmt(statement parser.Print) {
  value := i.evaluate(statement.Expression)

  if value == nil {
    fmt.Println("nil")
    return
  }

  fmt.Println(value)
}

func (i *Interpreter) executeVarStmt(statement parser.Var) {
  var value any = nil
  if statement.Initializer != nil {
    value = i.evaluate(statement.Initializer)
  }
  
  i.env.define(statement.Name.Lexeme, value)
}

func (i *Interpreter) executeBlockStmt(statement parser.Block, env *Environment) {
  previous := i.env
  i.env = env

  defer func() {
    i.env = previous
  }()
  
  for _, stmt := range statement.Statements {
    i.execute(stmt)
  } 
}

func (i *Interpreter) executeIfStmt(statement parser.If) {
  if isTruthy(i.evaluate(statement.Condition)) {
    i.execute(statement.ThenBranch)
  } else if statement.ElseBranch != nil {
    i.execute(statement.ElseBranch)
  }
}

func (i *Interpreter) executeWhileStmt(statement parser.While) {
  for isTruthy(i.evaluate(statement.Condition)) {
    i.execute(statement.Body)
  }
}
