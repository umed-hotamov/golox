package lexer

import "fmt"

type TokenType int

const (
  LEFT_PAREN TokenType = iota 
  RIGHT_PAREN
  LEFT_BRACE
  RIGHT_BRACE
  COMMA
  DOT
  MINUS
  PLUS
  SEMICOLON
  SLASH
  STAR


  IDENTIFIER
  STRING
  NUMBER

  AND
  CLASS
  ELSE
  FALSE
  FUN
  FOR
  IF
  NIL
  OR
  PRINT
  RETURN
  SUPER
  THIS
  TRUE
  VAR
  WHILE

  EOF
)

type Token struct {
  tokenType TokenType
  lexeme    string
  literal   any
  line      int
}

func NewToken(tokenType TokenType, lexeme string, literal any, line int) *Token {
  return &Token{
    tokenType: tokenType,
    lexeme: lexeme,
    literal: literal,
    line: line,
  }
}

func (tok *Token) String() string {
  return fmt.Sprintf("%v %s %v", tok.tokenType, tok.lexeme, tok.literal) 
}
