package ast 

import (
	"fmt"

	"github.com/umed-hotamov/golox/internal/lexer"
)

type Stmt interface {
  Ast
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

func (e Expression) Printer() string {
	return e.Expression.Printer() + ";"
}

func (p Print) Printer() string {
	return fmt.Sprintf("print %v", p.Expression.Printer())
}

func (v Var) Printer() string {
	return fmt.Sprintf("var %v = %v;", v.Name, v.Initializer.Printer())
}

func (b Block) Printer() string {
	var str string

	var i int
	for ; i < len(b.Statements)-1; i += 1 {
		str += b.Statements[i].Printer() + "\n"
	}
	str += b.Statements[i].Printer()

	return str
}

func (i If) Printer() string {
	return ""
}

func (w While) Printer() string {
	return ""
}

func (f Function) Printer() string {
	return fmt.Sprintf("fun %v", f.Name.Lexeme)
}

func (r Return) Printer() string {
	return fmt.Sprintf("return %v", r.Value.Printer())
}
