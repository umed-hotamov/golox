package interpreter

import "fmt"

type LoxClass struct {
  name string
}

func NewLoxClass(name string) *LoxClass {
  return &LoxClass{
    name: name,
  }
}

func (l *LoxClass) String() string {
  return fmt.Sprintf("class <%s>", l.name)
}
