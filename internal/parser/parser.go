package parser

import (
	"errors"
	"fmt"

	"github.com/umed-hotamov/golox/internal/lexer"
)

type Parser struct {
  tokens     []*lexer.Token
  statements []Stmt

  current    int
  HasError   bool
}

func NewParser(tokens []*lexer.Token) *Parser {
  return &Parser{
    tokens: tokens,
  }
}

func (p *Parser) Parse() []Stmt {
  for !p.eof() {
    p.statements = append(p.statements, p.declaration())
  }
  
  return p.statements
}

func (p *Parser) declaration() Stmt {
  defer p.errorRecovery()
  
  if p.match(lexer.VAR) {
    return p.varDeclaration()
  }

  return p.statement()
}

func (p *Parser) varDeclaration() Stmt {
  name := p.acceptToken(lexer.IDENTIFIER, "Expect variable name")
  
  var initializer Expr
  if p.match(lexer.EQUAL) {
    initializer = p.expression()
  }
  p.acceptToken(lexer.SEMICOLON, "Expect ; after variable declaration")

  return Var{*name, initializer}
}

func (p *Parser) statement() Stmt {
  if p.match(lexer.PRINT)      { return p.printStatement() }
  if p.match(lexer.LEFT_BRACE) { return p.block() }
  if p.match(lexer.IF)         { return p.ifStatement() }
  if p.match(lexer.WHILE)      { return p.whileStatement() }
  if p.match(lexer.FOR)        { return p.forStatement() }

  return p.expressionStatement()
}

func (p *Parser) printStatement() Stmt {
  value := p.expression()
  p.acceptToken(lexer.SEMICOLON, "Expect ; after value")

  return Print{value}
}

func (p *Parser) block() Block {
  var statements []Stmt

  for !p.eof() && !p.check(lexer.RIGHT_BRACE) {
    statements = append(statements, p.declaration())
  }

  p.acceptToken(lexer.RIGHT_BRACE, "Expect '}' after block")

  return Block{statements} 
}

func (p *Parser) expressionStatement() Stmt {
  expr := p.expression()
  p.acceptToken(lexer.SEMICOLON, "Expect ; after expression")

  return Expression{expr}
}

func (p *Parser) ifStatement() Stmt {
  p.acceptToken(lexer.LEFT_PAREN, "Expect ( after 'if')")
  condition := p.expression()
  p.acceptToken(lexer.RIGHT_PAREN, "Expect ) after if condition)")

  thenBranch := p.statement()
  var elseBranch Stmt
  if p.match(lexer.ELSE) {
    elseBranch = p.statement()
  }

  return If{condition, thenBranch, elseBranch}
} 

func (p *Parser) whileStatement() Stmt {
  p.acceptToken(lexer.LEFT_PAREN, "Expect ( after 'while'")
  condition := p.expression()
  p.acceptToken(lexer.RIGHT_PAREN, "Expect ) after condition")
  
  body := p.statement()

  return While{condition, body}
}

func (p *Parser) forStatement() Stmt {
  p.acceptToken(lexer.LEFT_PAREN, "Expect ( after 'for'")

  var initializer Stmt
  if p.match(lexer.SEMICOLON) {
    initializer = nil
  } else if p.match(lexer.VAR) {
    initializer = p.varDeclaration()
  } else if p.match(lexer.EQUAL) {
    initializer = p.expressionStatement()
  }
  
  var condition Expr
  if !p.check(lexer.SEMICOLON) {
    condition = p.expression()
  }
  p.acceptToken(lexer.SEMICOLON, "Expect ; after loop condition")

  var increment Expr
  if !p.check(lexer.RIGHT_PAREN) {
    increment = p.expression()
  }
  p.acceptToken(lexer.RIGHT_PAREN, "Expect ) after for clauses")
  
  body := p.statement()
  if increment != nil {
    body = Block{[]Stmt{body, Expression{increment}}}
  }
  if condition == nil {
    condition = Literal{true}
  }
  body = While{condition, body}

  if initializer != nil {
    body = Block{[]Stmt{initializer, body}}
  }

  return body
}

func (p *Parser) expression() Expr {
  return p.assignment()
}

func (p *Parser) assignment() Expr {
  expr := p.or()

  if p.match(lexer.EQUAL) {
    equals := p.previous()
    value := p.equality()

    switch expr.(type) {
    case Variable:
      name := expr.(Variable).Name
      return Assign{name, value}
    }

    p.error(equals, errors.New("Invalid assignment target"))
  }

  return expr
}

func (p *Parser) or() Expr {
  expr := p.and()

  for p.match(lexer.OR) {
    operator := p.previous()
    right := p.and()
    expr = Logical{expr, *operator, right}
  }
  
  return expr
}

func (p *Parser) and() Expr {
  expr := p.equality()

  for p.match(lexer.AND) {
    operator := p.previous()
    right := p.equality()
    expr = Logical{expr, *operator, right}
  }

  return expr
}

func (p *Parser) equality() Expr {
  expr := p.comprasion()

  for p.match(lexer.BANG_EQUAL, lexer.EQUAL_EQUAL) {
    operator := p.previous()
    right := p.comprasion()

    expr = Binary{expr, *operator, right}
  }

  return expr
}

func (p *Parser) comprasion() Expr {
  expr := p.term()

  for p.match(lexer.GREATER, lexer.GREATER_EQUAL, lexer.LESS, lexer.LESS_EQUAL) {
    operator := p.previous()
    right := p.term()

    expr = Binary{expr, *operator, right}
  }

  return expr
}

func (p *Parser) term() Expr {
  expr := p.factor()

  for p.match(lexer.MINUS, lexer.PLUS) {
    operator := p.previous()
    right := p.factor()

    expr = Binary{expr, *operator, right}
  }

  return expr
}

func (p *Parser) factor() Expr {
  expr := p.unary()

  for p.match(lexer.STAR, lexer.SLASH) {
    operator := p.previous()
    right := p.unary()

    expr = Binary{expr, *operator, right}
  }

  return expr
}

func (p *Parser) unary() Expr {
  if p.match(lexer.BANG, lexer.MINUS) {
    operator := p.previous()
    right := p.unary()

    return Unary{*operator, right}
  }

  return p.primary()
}

func (p *Parser) primary() Expr {
  if p.match(lexer.TRUE)  { return Literal{true}  }
  if p.match(lexer.FALSE) { return Literal{false} }
  if p.match(lexer.NIL)   { return Literal{nil}   }
  
  if p.match(lexer.NUMBER, lexer.STRING) {
    return Literal{p.previous().Literal}
  }
  if p.match(lexer.IDENTIFIER) {
    return Variable{*p.previous()}
  }

  if p.match(lexer.LEFT_PAREN) {
    expr := p.expression()
    p.acceptToken(lexer.RIGHT_PAREN, "Expect ')' after expression")
    return Grouping{expr}
  }
  
  p.parseError("Expect expression")
  return nil
}

func (p *Parser) errorRecovery() {
  if err := recover(); err != nil {
    p.error(p.peek(), fmt.Errorf("%v", err))
    p.synchronize()
  }
}

func (p *Parser) synchronize() {
  p.advance()

  for !p.eof() {
    if p.previous().TokenType == lexer.SEMICOLON {
      return
    }
    switch p.peek().TokenType {
    case lexer.IF, lexer.VAR, lexer.FUN, lexer.FOR,
         lexer.WHILE, lexer.PRINT, lexer.CLASS, lexer.RETURN:
    }
    p.advance()
  }
}
