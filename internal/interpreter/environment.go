package interpreter

import (
	"fmt"

	"github.com/umed-hotamov/golox/internal/lexer"
)

type Environment struct {
  objects   map[string]any
  enclosing *Environment
}

func NewEnvironment() *Environment {
  return &Environment{
    objects: make(map[string]any),
  }
}

func NewEnclosingEnvironment(env *Environment) *Environment {
  enclosingEnv := NewEnvironment()
  enclosingEnv.enclosing = env
  return enclosingEnv
} 

func (e *Environment) define(name string, value any) {
  e.objects[name] = value
}

func (e *Environment) get(token lexer.Token) any {
  if value, ok := e.objects[token.Lexeme]; ok {
    return value
  }

  if e.enclosing != nil {
    return e.enclosing.get(token)
  }

  panic(fmt.Sprintf("Undefined variable %s\n", token.Lexeme))
}

func (e *Environment) getAt(distance int, value string) any {
  ancestor := e.ancestor(distance)
  return ancestor.objects[value]
}

func (e *Environment) assignAt(distance int, name lexer.Token, value any) {
  ancestor := e.ancestor(distance)
  ancestor.objects[name.Lexeme] = value
}

func (e *Environment) ancestor(distance int) *Environment {
  a := e
  for range distance {
    a = a.enclosing
  }

  return a
}

func (e *Environment) assign(name lexer.Token, value any) {
  if _, ok := e.objects[name.Lexeme]; ok {
    e.objects[name.Lexeme] = value
    return
  }

  if e.enclosing != nil {
    e.enclosing.assign(name, value)
    return
  }

  panic(fmt.Sprintf("Undefined variable %s\n", name.Lexeme))
}
