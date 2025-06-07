package lexer

import (
	"fmt"
	"strconv"
	"strings"
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
  
  l.tokens = append(l.tokens, NewToken(EOF, "", nil, l.line))

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
      l.startColumn = l.current
      l.line += 1
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

func (l *Lexer) isAlpha(c byte) bool {
  return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' 
}

func (l *Lexer) isDigit(c byte) bool {
  return c >= '0' && c <= '9'
}

func (l *Lexer) isAlphaNumeric(c byte) bool {
  return l.isAlpha(c) || l.isDigit(c)
}

func (l *Lexer) acceptBlockComments() {
  depth := 1
  for !l.eof() && depth > 0 {
    if l.peek() == '*' && l.peekNext() == '/' {
      depth -= 1
      l.advance()
    } else if l.peek() == '/' && l.peekNext() == '*' {
      depth += 1
      l.advance()
    } else if l.peek() == '\n' {
      l.line += 1
      l.startColumn = l.current
    }

    l.advance()
  }

  if depth > 0 {
    l.error("Unterminated block comment")
  }
}

func (l *Lexer) acceptIdentifier() {
  for !l.eof() && l.isAlphaNumeric(l.peek()) {
    l.advance()
  }

  literal := l.source[l.start:l.current]
  
  tokenType, ok := keywords[literal]
  if !ok {
    tokenType = IDENTIFIER
  }

  l.addTokenLiteral(tokenType, literal)
}

func (l *Lexer) acceptNumber() {
  digits := "0123456789"
  l.acceptRun(digits)
  if l.accept('.') && l.isDigit(l.peek()) {
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
      l.line += 1
      l.startColumn = l.current
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

func (l *Lexer) error(message string) {
  fmt.Printf("[line: %d, column: %d] Error: %s\n", l.lineStart, l.current - l.startColumn, message)
  fmt.Printf("|   %c\n", l.source[l.start])
  fmt.Printf("|---^\n")
  l.HasError = true
}
