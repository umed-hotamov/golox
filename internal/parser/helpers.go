package parser

import (
	"fmt"

	"github.com/umed-hotamov/golox/internal/lexer"
)


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

func (p *Parser) acceptToken(tokenType lexer.TokenType, message string) *lexer.Token {
  if p.check(tokenType) {
    return p.advance()
  }

  p.parseError(message)
  return nil
}

func (p *Parser) parseError(message string) {
  p.HasError = true
  panic(message)
}

func (p *Parser) error(token *lexer.Token, err error) {
  fmt.Printf("[line: %d] Error: %s\n", token.Line, err.Error())
}
