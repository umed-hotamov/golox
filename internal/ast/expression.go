package ast

import (
	"fmt"

	"github.com/umed-hotamov/golox/internal/lexer"
)

type Expr interface {
  Ast
}

type Binary struct {
	Left     Expr
	Operator lexer.Token
	Right    Expr
}

type Grouping struct {
	Expr Expr
}

type Literal struct {
	Value any
}

type Unary struct {
	Operator lexer.Token
	Right    Expr
}

type Variable struct {
	Name lexer.Token
}

type Assign struct {
	Name  lexer.Token
	Value Expr
}

type Logical struct {
	Left     Expr
	Operator lexer.Token
	Right    Expr
}

type Call struct {
	Callee    Expr
	Paren     lexer.Token
	Arguments []Expr
}

func (b Binary) Printer() string {
	return fmt.Sprintf("(%v %v %v)", b.Operator.Lexeme, b.Left.Printer(), b.Right.Printer())
}

func (g Grouping) Printer() string {
	return fmt.Sprintf("(group %v)", g.Expr.Printer())
}

func (l Literal) Printer() string {
	if l.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", l.Value)
}

func (u Unary) Printer() string {
	return fmt.Sprintf("(%v %v)", u.Operator.Lexeme, u.Right.Printer())
}

func (v Variable) Printer() string {
	return fmt.Sprint(v.Name.Lexeme)
}

func (a Assign) Printer() string {
	return fmt.Sprintf("(%v %v)", a.Value.Printer(), a.Name)
}

func (l Logical) Printer() string {
	return fmt.Sprintf("%v %v %v", l.Left.Printer(), l.Operator.Lexeme, l.Right.Printer())
}

func (c Call) Printer() string {
	s := fmt.Sprintf("%v", c.Callee.Printer())
	s += "("
	for _, a := range c.Arguments {
		s += a.Printer()
	}
	s += ")"

	return s
}
