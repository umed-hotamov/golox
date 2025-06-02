package parser

import (
	"fmt"

	"github.com/umed-hotamov/golox/internal/lexer"
)

type Stmt interface {
	String() string
}

type Expression struct {
	Expression Expr
}

type Print struct {
	Expression Expr
}

type Var struct {
	Name        lexer.Token
	Initializer Expr
}

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

type Block struct {
	Statements []Stmt
}

type While struct {
	Condition Expr
	Body      Stmt
}

type Function struct {
	Name   lexer.Token
	Params []lexer.Token
	Body   Block
}

type Return struct {
	Keyword lexer.Token
	Value   Expr
}

func (e Expression) String() string {
	return e.Expression.String() + ";"
}

func (p Print) String() string {
	return fmt.Sprintf("print %v", p.Expression.String())
}

func (v Var) String() string {
	return fmt.Sprintf("var %v = %v;", v.Name, v.Initializer.String())
}

func (b Block) String() string {
	var str string

	var i int
	for ; i < len(b.Statements)-1; i += 1 {
		str += b.Statements[i].String() + "\n"
	}
	str += b.Statements[i].String()

	return str
}

func (i If) String() string {
	return ""
}

func (w While) String() string {
	return ""
}

func (f Function) String() string {
	return fmt.Sprintf("fun %v", f.Name.Lexeme)
}

func (r Return) String() string {
	return fmt.Sprintf("return %v", r.Value.String())
}
