package lexer

import (
  "strconv"
	"strings"
)

func (l *Lexer) isAlpha(c byte) bool {
  return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' 
}

func (l *Lexer) isDigit(c byte) bool {
  return c >= '0' && c <= '9'
}

func (l *Lexer) isAlphaNumeric(c byte) bool {
  return l.isAlpha(c) || l.isDigit(c)
}

func (l *Lexer) nextLine() {
  l.startColumn = l.current
  l.line += 1
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
      l.nextLine()
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
      l.nextLine()
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
  column := l.current - l.startColumn
  l.tokens = append(l.tokens, NewToken(tokenType, lexeme, literal, l.line, column))
}
