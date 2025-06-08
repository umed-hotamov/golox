package lexer

import (
	"fmt"
)

type Lexer struct {
  source      string
  tokens      []*Token

  line        int
  lineStart   int
  
  start       int
  startColumn int
  current     int

  HasError    bool
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
    l.lineStart = l.line
    l.fetchToken()
  }
  
  l.tokens = append(l.tokens, NewToken(EOF, "", nil, l.line, 0))

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
      } else if l.accept('*') {
        l.acceptBlockComments()
      } else {
        l.addToken(SLASH)
      }
    case '"':
      l.acceptString()
    case '\n':
      l.nextLine()
    case ' ', '\r', '\t':
    default:
      if l.isDigit(c) {
        l.acceptNumber()
      } else if l.isAlpha(c) {
        l.acceptIdentifier()
      } else {
        l.error("Unexpected character")
      }
  }
}

func (l *Lexer) error(message string) {
  fmt.Printf("[line: %d, column: %d] Error: %s\n", l.lineStart, l.current - l.startColumn, message)
  fmt.Printf("|   %c\n", l.source[l.start])
  fmt.Printf("|---^\n")
  l.HasError = true
}
