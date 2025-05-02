package lexer

import "fmt"

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

func (l *Lexer) Lex() []*Token {
  for !l.eof() {

    l.start = l.current
    l.fetchToken()
  }

  return l.tokens
}

func (l *Lexer) fetchToken() {
  c := l.source[l.current]
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
    case '/':
      l.addToken(SLASH)
    case ';':
      l.addToken(SEMICOLON)
    default:
      l.error("Unexpected character")
  }
}

func (l *Lexer) advance() byte {
  current := l.current
  l.current += 1
  return l.source[current]
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
  fmt.Printf("[line: %d] Error: %s", l.line, message)
}
