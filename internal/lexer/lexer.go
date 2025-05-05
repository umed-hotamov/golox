package lexer

import (
	"fmt"
	"strconv"
	"strings"
)

type Lexer struct {
  source   string
  tokens   []*Token

  line     int
  start    int
  current  int

  hasError bool
}

func NewLexer(source string) *Lexer {
  tokens := make([]*Token, 0)

  return &Lexer{
    source: source,
    tokens: tokens,
    line: 1,
  }
}

var keywords = map[string]TokenType{
  "and":    AND,
  "or":     OR,
  "class":  CLASS,
  "else":   ELSE,
  "false":  FALSE,
  "true":   TRUE,
  "if":     IF,
  "nil":    NIL,
  "for":    FOR,
  "fun":    FUN,
  "print":  PRINT,
  "return": RETURN,
  "super":  SUPER,
  "this":   THIS,
  "var":    VAR,
  "while":  WHILE,
}

func (l *Lexer) Lex() []*Token {
  for !l.eof() {

    l.start = l.current
    l.fetchToken()
  }

  return l.tokens
}

func (l *Lexer) fetchToken() {
  c := l.advance()
  switch c {
    case '{':
      l.addToken(LEFT_BRACE)
    case '}':
      l.addToken(RIGHT_BRACE)
    case '(':
      l.addToken(LEFT_PAREN)
    case ')':
      l.addToken(RIGHT_PAREN)
    case ',':
      l.addToken(COMMA)
    case '.':
      l.addToken(DOT)
    case '+':
      l.addToken(PLUS)
    case '-':
      l.addToken(MINUS)
    case '*':
      l.addToken(STAR)
    case ';':
      l.addToken(SEMICOLON)
    case '!':
      if l.accept('=') {
        l.addToken(BANG_EQUAL)
      } else {
        l.addToken(BANG)
      }
    case '=':
      if l.accept('=') {
        l.addToken(EQUAL_EQUAL)
      } else {
        l.addToken(EQUAL)
      }
    case '>':
      if l.accept('=') {
        l.addToken(GREATER_EQUAL)
      } else {
        l.addToken(GREATER)
      }
    case '<':
      if l.accept('=') {
        l.addToken(LESS_EQUAL)
      } else {
        l.addToken(LESS)
      }
    case '/':
      if l.accept('/') {
        l.skipTo('\n')
      } else {
        l.addToken(SLASH)
      }
    case '"':
      l.acceptString()
    case '\n':
      l.incLine()
    case ' ', '\r', '\t':
    default:
      if l.isDigit(c) {
        l.acceptNumber()
      } else {
        l.error("Unexpected character")
      }
  }
}

func (l *Lexer) isDigit(c byte) bool {
  return c >= '0' && c <= '9'
}

func (l *Lexer) acceptNumber() {
  digits := "0123456789"
  l.acceptRun(digits)
  if l.accept('.') && !l.isDigit(l.peekNext()) {
    l.acceptRun(digits)
  }

  number, _ := strconv.ParseFloat(l.source[l.start:l.current], 64)
  l.addTokenLiteral(NUMBER, number)
}

func (l *Lexer) acceptString() {
  l.skipTo('"')
  if l.eof() {
    l.error("Unterminated string")
    return
  }

  l.advance()
  l.addTokenLiteral(STRING, l.source[l.start + 1:l.current-1])
}

func (l *Lexer) accept(valid byte) bool {
  if l.peek() == valid {
    l.advance()
    return true
  }

  return false
}

func (l *Lexer) acceptRun(valid string) {
  for !l.eof() && strings.ContainsRune(valid, rune(l.peek())) {
    l.advance()
  }
}

func (l *Lexer) skipTo(to byte) {
  for !l.eof() && l.peek() != to {
    if l.peek() == '\n' {
      l.incLine()
    }
    l.advance()
  }
}

func (l *Lexer) advance() byte {
  current := l.current
  l.current += 1
  return l.source[current]
}

func (l *Lexer) peek() byte {
  if l.eof() {
    return 0 
  } 

  return l.source[l.current]
}

func (l *Lexer) peekNext() byte {
  if l.current + 1 >= len(l.source) {
    return 0
  }
  
  return l.source[l.current + 1]
}

func (l *Lexer) eof() bool {
  return l.current >= len(l.source)
}

func (l *Lexer) addToken(tokenType TokenType) {
  l.addTokenLiteral(tokenType, nil)
}

func (l *Lexer) addTokenLiteral(tokenType TokenType, literal any) {
  lexeme := l.source[l.start:l.current]
  l.tokens = append(l.tokens, NewToken(tokenType, lexeme, literal, l.line))
}

func (l *Lexer) incLine() {
  l.line += 1
}

func (l *Lexer) error(message string) {
  fmt.Printf("[line: %d] Error: %s\n", l.line, message)
}
