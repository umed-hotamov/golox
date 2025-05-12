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

func (b *Binary) String() string {
  return fmt.Sprintf("(%v %v %v)", b.Operator.Lexeme, b.Left.String(), b.Right.String())
}

func (g *Grouping) String() string {
  return fmt.Sprintf("(group %v)", g.Expr.String())
}

func (l *Literal) String() string {
  if l.Value == nil {
    return "nil"
  }
  return fmt.Sprintf("%v", l.Value)
}

func (u *Unary) String() string {
  return fmt.Sprintf("(%v %v)", u.Operator.Lexeme, u.Right.String())
}
