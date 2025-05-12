package parser

import (
	"fmt"

	"github.com/umed-hotamov/golox/internal/lexer"
)

type Parser struct {
  tokens []*lexer.Token

  current int
}

func NewParser(tokens []*lexer.Token) *Parser {
  return &Parser{
    tokens: tokens,
  }
}

func (p *Parser) Parse() Expr {
  defer p.errorRecovery()
  return p.expression()
}

func (p *Parser) peek() *lexer.Token {
  return p.tokens[p.current]
}

func (p *Parser) eof() bool {
  return p.peek().TokenType == lexer.EOF
}

func (p *Parser) previous() *lexer.Token {
  return p.tokens[p.current - 1]
}

func (p *Parser) advance() *lexer.Token {
  if !p.eof() {
    p.current += 1
  }

  return p.previous()
}

func (p *Parser) check(tokenType lexer.TokenType) bool {
  if p.eof() {
    return false
  }

  return tokenType == p.peek().TokenType
}

func (p *Parser) match(types... lexer.TokenType) bool {
  for _, tokenType := range types {
    if p.check(tokenType) {
      p.advance()
      return true
    }
  }

  return false
}

func (p *Parser) expression() Expr {
  return p.equality()
}

func (p *Parser) equality() Expr {
  expr := p.comprasion()

  for p.match(lexer.BANG_EQUAL, lexer.EQUAL_EQUAL) {
    operator := p.previous()
    right := p.comprasion()

    expr = &Binary{expr, *operator, right}
  }

  return expr
}

func (p *Parser) comprasion() Expr {
  expr := p.term()

  for p.match(lexer.GREATER, lexer.GREATER_EQUAL, lexer.LESS, lexer.LESS_EQUAL) {
    operator := p.previous()
    right := p.term()

    expr = &Binary{expr, *operator, right}
  }

  return expr
}

func (p *Parser) term() Expr {
  expr := p.factor()

  for p.match(lexer.MINUS, lexer.PLUS) {
    operator := p.previous()
    right := p.factor()

    expr = &Binary{expr, *operator, right}
  }

  return expr
}

func (p *Parser) factor() Expr {
  expr := p.unary()

  for p.match(lexer.STAR, lexer.SLASH) {
    operator := p.previous()
    right := p.unary()

    expr = &Binary{expr, *operator, right}
  }

  return expr
}

func (p *Parser) unary() Expr {
  if p.match(lexer.BANG, lexer.MINUS) {
    operator := p.previous()
    right := p.unary()

    return &Unary{*operator, right}
  }

  return p.primary()
}

func (p *Parser) primary() Expr {
  if p.match(lexer.TRUE)  { return &Literal{true}  }
  if p.match(lexer.FALSE) { return &Literal{false} }
  if p.match(lexer.NIL)   { return &Literal{nil}   }
  
  if p.match(lexer.NUMBER, lexer.STRING) {
    return &Literal{p.previous().Literal}
  }

  if p.match(lexer.LEFT_PAREN) {
    expr := p.expression()
    p.acceptToken(lexer.RIGHT_PAREN, "Expect ')' after expression")
    return &Grouping{expr}
  }
  
  p.parseError(p.peek(), "Expect expression")
  return nil
}

func (p *Parser) acceptToken(tokenType lexer.TokenType, message string) {
  if p.check(tokenType) {
    p.advance()
    return
  }

  p.parseError(p.peek(), message)
}

func (p *Parser) parseError(token *lexer.Token, message string) {
  panic(message)
}

func (p *Parser) errorRecovery() {
  if err := recover(); err != nil {
    p.error(p.peek(), fmt.Errorf("%v", err))
  }
}

func (p *Parser) error(token *lexer.Token, err error) {
  fmt.Printf("[line: %d] Error: %s\n", token.Line, err.Error())
}
