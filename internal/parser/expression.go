package parser

import (
	"fmt"

	"github.com/umed-hotamov/golox/internal/lexer"
)

type Expr interface {
	String() string
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

func (b Binary) String() string {
	return fmt.Sprintf("(%v %v %v)", b.Operator.Lexeme, b.Left.String(), b.Right.String())
}

func (g Grouping) String() string {
	return fmt.Sprintf("(group %v)", g.Expr.String())
}

func (l Literal) String() string {
	if l.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", l.Value)
}

func (u Unary) String() string {
	return fmt.Sprintf("(%v %v)", u.Operator.Lexeme, u.Right.String())
}

func (v Variable) String() string {
	return fmt.Sprint(v.Name.Lexeme)
}

func (a Assign) String() string {
	return fmt.Sprintf("(%v %v)", a.Value.String(), a.Name)
}

func (l Logical) String() string {
	return fmt.Sprintf("%v %v %v", l.Left.String(), l.Operator.Lexeme, l.Right.String())
}

func (c Call) String() string {
	s := fmt.Sprintf("%v", c.Callee.String())
	s += "("
	for _, a := range c.Arguments {
		s += a.String()
	}
	s += ")"

	return s
}
