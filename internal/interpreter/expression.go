package interpreter

import (
	"fmt"

	"github.com/umed-hotamov/golox/internal/ast"
	"github.com/umed-hotamov/golox/internal/lexer"
)

func (i *Interpreter) evaluate(expression ast.Expr) any {
	switch expression.(type) {
	case ast.Literal:
		return i.evaluateLiteral(expression.(ast.Literal))
	case ast.Grouping:
		return i.evaluateGrouping(expression.(ast.Grouping))
	case ast.Unary:
		return i.evaluateUnary(expression.(ast.Unary))
	case ast.Binary:
		return i.evaluateBinary(expression.(ast.Binary))
	case ast.Variable:
		return i.evaluateVariable(expression.(ast.Variable))
	case ast.Assign:
		return i.evaluateAssign(expression.(ast.Assign))
	case ast.Logical:
		return i.evaluateLogical(expression.(ast.Logical))
	case ast.Call:
		return i.evaluateCall(expression.(ast.Call))
	}

	return nil
}

func (i *Interpreter) evaluateLiteral(expression ast.Literal) any {
	return expression.Value
}

func (i *Interpreter) evaluateGrouping(expression ast.Grouping) any {
	return i.evaluate(expression.Expr)
}

func (i *Interpreter) evaluateUnary(expression ast.Unary) any {
	right := i.evaluate(expression.Right)

	switch expression.Operator.TokenType {
	case lexer.BANG:
		return !isTruthy(right)
	case lexer.MINUS:
		return -right.(float64)
	}

	return nil
}

func (i *Interpreter) evaluateBinary(expression ast.Binary) any {
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

func (i *Interpreter) evaluateVariable(expression ast.Variable) any {
	return i.lookUpVariable(expression.Name, expression)
}

func (i *Interpreter) lookUpVariable(name lexer.Token, expression ast.Expr) any {
	distance, ok := i.locals[expression]
	if ok {
		return i.env.getAt(distance, name.Lexeme)
	}

	return i.globals.get(name)
}

func (i *Interpreter) evaluateAssign(expression ast.Assign) any {
	value := i.evaluate(expression.Value)

	distance, ok := i.locals[expression]
	if ok {
		i.env.assignAt(distance, expression.Name, value)
	} else {
		i.globals.assign(expression.Name, value)
	}

	return value
}

func (i *Interpreter) evaluateLogical(expression ast.Logical) any {
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

func (i *Interpreter) evaluateCall(expression ast.Call) any {
	callee := i.evaluate(expression.Callee)

	var arguments []any
	for _, arg := range expression.Arguments {
		arguments = append(arguments, i.evaluate(arg))
	}
	_, ok := callee.(Callable)
	if !ok {
		runtimeError(expression.Paren, "Call only call functions and classes")
	}
	function := callee.(Callable)
	if function.arity() != len(arguments) {
		runtimeError(expression.Paren, fmt.Sprintf("Expected %d, arguments got %d", function.arity(), len(arguments)))
	}

	return function.call(i, arguments)
}
