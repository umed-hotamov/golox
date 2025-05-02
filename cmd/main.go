package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
  args := os.Args[1:]

  if len(args) > 1 {
    fmt.Println("Usage: glox [source]")
  }

  if len(args) == 1 {
    runFile(args[1])
  } else {
    runPrompt()
  }
}

func runFile(filename string) {

}

func runPrompt() {
  scanner := bufio.NewScanner(os.Stdin)
  
  for {
    fmt.Print("~> ")
    
    line := scanner.Scan()
    if !line {
      return
    }

    run()
  }
}

func run() {

}
