package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/umed-hotamov/golox/internal/interpreter"
	"github.com/umed-hotamov/golox/internal/lexer"
	"github.com/umed-hotamov/golox/internal/parser"
)

func main() {
  args := os.Args[1:]

  if len(args) > 1 {
    fmt.Println("Usage: glox [source]")
  }
  
  if len(args) == 1 {
    runFile(args[0])
  } else {
    runPrompt()
  }
}

func runFile(filename string) {
  data, err := os.ReadFile(filename)
  if err != nil {
    log.Fatalf("failed to read file: %s", err)
  }
  source := string(data)

  interpreter := interpreter.NewInterpreter()
  run(source, interpreter)
}

func runPrompt() {
  scanner := bufio.NewScanner(os.Stdin)
  interpreter := interpreter.NewInterpreter()
  
  for {
    fmt.Print("golox~~>  ")
    
    isNotEnd := scanner.Scan()
    if !isNotEnd {
      return
    }
    line := scanner.Text()
    if line == "exit" {
      break
    }

    run(line, interpreter)
  }
}

func run(source string, interpreter *interpreter.Interpreter) {
  lexer := lexer.NewLexer(source)
  tokens := lexer.Lex()

  parser := parser.NewParser(tokens)
  statements := parser.Parse()
  if lexer.HasError {
    return
  }
  if parser.HasError {
    return
  }

  interpreter.Interpret(statements)
}
