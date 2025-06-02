package interpreter

import (
	"fmt"

	"github.com/umed-hotamov/golox/internal/lexer"
	"github.com/umed-hotamov/golox/internal/parser"
)

func (i *Interpreter) evaluate(expression parser.Expr) any {
  switch expression.(type) {
  case parser.Literal:
    return i.evaluateLiteral(expression.(parser.Literal))
  case parser.Grouping:
    return i.evaluateGrouping(expression.(parser.Grouping))
  case parser.Unary:
    return i.evaluateUnary(expression.(parser.Unary))
  case parser.Binary:
    return i.evaluateBinary(expression.(parser.Binary))
  case parser.Variable:
    return i.evaluateVariable(expression.(parser.Variable))
  case parser.Assign:
    return i.evaluateAssign(expression.(parser.Assign))
  case parser.Logical:
    return i.evaluateLogical(expression.(parser.Logical))
  case parser.Call:
    return i.evaluateCall(expression.(parser.Call))
  }

  return nil
}

func (i *Interpreter) evaluateLiteral(expression parser.Literal) any {
  return expression.Value
}

func (i *Interpreter) evaluateGrouping(expression parser.Grouping) any {
  return i.evaluate(expression.Expr)
}

func (i *Interpreter) evaluateUnary(expression parser.Unary) any {
  right := i.evaluate(expression.Right)

  switch expression.Operator.TokenType {
  case lexer.BANG:
    return !isTruthy(right)
  case lexer.MINUS:
    return -right.(float64)
  }

  return nil
}

func (i *Interpreter) evaluateBinary(expression parser.Binary) any {
  left := i.evaluate(expression.Left)
  right := i.evaluate(expression.Right)

  switch expression.Operator.TokenType {
  case lexer.EQUAL_EQUAL:
    return isEqual(left, right)
  case lexer.BANG_EQUAL:
    return !isEqual(left, right)
  case lexer.GREATER:
    return left.(float64) > right.(float64)
  case lexer.LESS:
    return left.(float64) < right.(float64)
  case lexer.GREATER_EQUAL:
    return left.(float64) >= right.(float64)
  case lexer.LESS_EQUAL:
    return left.(float64) <= right.(float64)
  case lexer.STAR:
    return left.(float64) * right.(float64)
  case lexer.SLASH:
    return left.(float64) / right.(float64)
  case lexer.MINUS:
    return left.(float64) - right.(float64)
  case lexer.PLUS:
    if isNumber(left) && isNumber(right) {
      return left.(float64) + right.(float64)
    }
    if isString(left) && isString(right) {
      return left.(string) + right.(string)
    }

    runtimeError(expression.Operator, "Operands must be either numbers or strings")
  }

  return nil
}

func (i *Interpreter) evaluateVariable(expression parser.Variable) any {
  return i.env.get(expression.Name)
}

func (i *Interpreter) evaluateAssign(expression parser.Assign) any {
  value := i.evaluate(expression.Value)
  i.env.assign(expression.Name, value)
  
  return value
}

func (i *Interpreter) evaluateLogical(expression parser.Logical) any {
  left := i.evaluate(expression.Left)
  if expression.Operator.TokenType == lexer.OR {
    if isTruthy(left) {
      return left
    }
  } else {
      if !isTruthy(left) {
        return left
      }
  }

  return i.evaluate(expression.Right)
}

func (i *Interpreter) evaluateCall(expression parser.Call) any {
  callee := i.evaluate(expression.Callee)

  var arguments []any
  for _, arg := range expression.Arguments {
    arguments = append(arguments, i.evaluate(arg))
  }
  _, ok := callee.(Callable)
  if !ok {
    runtimeError(expression.Paren, "Call only  call functions and classes")
  }
  function := callee.(Callable)
  if function.arity() != len(arguments) {
    runtimeError(expression.Paren, fmt.Sprintf("Expected %d, arguments got %d", function.arity(), len(arguments)))
  }

  return function.call(i, arguments)
}
