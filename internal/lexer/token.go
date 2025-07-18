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

  BANG
  EQUAL
  LESS
  GREATER
  BANG_EQUAL
  EQUAL_EQUAL
  LESS_EQUAL
  GREATER_EQUAL

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
  TokenType TokenType
  Lexeme    string
  Literal   any
  Line      int
  Column    int
}

func NewToken(tokenType TokenType, lexeme string, literal any, line int, column int) *Token {
  return &Token{
    TokenType: tokenType,
    Lexeme:    lexeme,
    Literal:   literal,
    Line:      line,
    Column:    column,
  }
}

func (tok *Token) String() string {
  return fmt.Sprintf("[type = %v lexeme = %s literal = %v line = %d column = %d]", tok.TokenType, tok.Lexeme, tok.Literal, tok.Line, tok.Column) 
}
