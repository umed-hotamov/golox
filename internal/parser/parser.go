package parser

import (
	"errors"
	"fmt"

	"github.com/umed-hotamov/golox/internal/lexer"
	"github.com/umed-hotamov/golox/internal/ast"
)

type Parser struct {
	tokens     []*lexer.Token
	statements []ast.Stmt

	current  int
	HasError bool
}

func NewParser(tokens []*lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parse() []ast.Stmt {
	for !p.eof() {
		p.statements = append(p.statements, p.declaration())
	}

	return p.statements
}

func (p *Parser) declaration() ast.Stmt {
	defer p.errorRecovery()

	if p.match(lexer.VAR) {
		return p.varDeclaration()
	}
	if p.match(lexer.FUN) {
		return p.function("function")
	}
  if p.match(lexer.CLASS) {
    return p.classDeclaration()
  }


	return p.statement()
}

func (p *Parser) varDeclaration() ast.Stmt {
	name := p.acceptToken(lexer.IDENTIFIER, "Expect variable name")

	var initializer ast.Expr
	if p.match(lexer.EQUAL) {
		initializer = p.expression()
	}
	p.acceptToken(lexer.SEMICOLON, "Expect ; after variable declaration")

  return ast.Var{Name: *name, Initializer: initializer}
}

func (p *Parser) function(kind string) ast.Stmt {
	name := p.acceptToken(lexer.IDENTIFIER, "Expect "+kind+" name")

	p.acceptToken(lexer.LEFT_PAREN, "Expect ( after "+kind+" name")
	var parameters []lexer.Token
	if !p.check(lexer.RIGHT_PAREN) {
		parameters = append(parameters, *p.acceptToken(lexer.IDENTIFIER, "Expect parameter name"))
	}

	for p.match(lexer.COMMA) {
		parameters = append(parameters, *p.acceptToken(lexer.IDENTIFIER, "Expect parameter name"))
		if len(parameters) > 255 {
			p.error(p.peek(), errors.New("Can't have more than 255 parameters"))
		}

		if p.check(lexer.RIGHT_PAREN) {
			break
		}
	}
	p.acceptToken(lexer.RIGHT_PAREN, "Expect ')' after arguments")

	p.acceptToken(lexer.LEFT_BRACE, "Expect '{' before "+kind+" body")
	body := p.block()

  return ast.Function{Name: *name, Params: parameters, Body: ast.Block{Statements: body.Statements}}
}

func (p *Parser) classDeclaration() ast.Stmt {
  name := p.acceptToken(lexer.IDENTIFIER, "Expect class name")  
  p.acceptToken(lexer.LEFT_BRACE, "Expect '{' before class body")

  var methods []ast.Function
  for !p.check(lexer.RIGHT_BRACE) && !p.eof() {
    methods = append(methods, p.function("method").(ast.Function))
  }
  p.acceptToken(lexer.RIGHT_BRACE, "Expect '}' after class body")
 
  return ast.Class{Name: *name, Methods: methods}
}

func (p *Parser) statement() ast.Stmt {
	if p.match(lexer.PRINT) {
		return p.printStatement()
	}
	if p.match(lexer.LEFT_BRACE) {
		return p.block()
	}
	if p.match(lexer.IF) {
		return p.ifStatement()
	}
	if p.match(lexer.WHILE) {
		return p.whileStatement()
	}
	if p.match(lexer.FOR) {
		return p.forStatement()
	}
	if p.match(lexer.RETURN) {
		return p.returnStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() ast.Stmt {
	value := p.expression()
	p.acceptToken(lexer.SEMICOLON, "Expect ; after value")

	return ast.Print{Expression: value}
}

func (p *Parser) block() ast.Block {
	var statements []ast.Stmt

	for !p.eof() && !p.check(lexer.RIGHT_BRACE) {
		statements = append(statements, p.declaration())
	}

	p.acceptToken(lexer.RIGHT_BRACE, "Expect '}' after block")

	return ast.Block{Statements: statements}
}

func (p *Parser) expressionStatement() ast.Stmt {
	expr := p.expression()
	p.acceptToken(lexer.SEMICOLON, "Expect ; after expression")

	return ast.Expression{Expression: expr}
}

func (p *Parser) ifStatement() ast.Stmt {
	p.acceptToken(lexer.LEFT_PAREN, "Expect ( after 'if'")
	condition := p.expression()
	p.acceptToken(lexer.RIGHT_PAREN, "Expect ) after if condition")

	thenBranch := p.statement()
	var elseBranch ast.Stmt
	if p.match(lexer.ELSE) {
		elseBranch = p.statement()
	}

  return ast.If{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}
}

func (p *Parser) whileStatement() ast.Stmt {
	p.acceptToken(lexer.LEFT_PAREN, "Expect ( after 'while'")
	condition := p.expression()
	p.acceptToken(lexer.RIGHT_PAREN, "Expect ) after condition")

	body := p.statement()

  return ast.While{Condition: condition, Body: body}
}

func (p *Parser) forStatement() ast.Stmt {
	p.acceptToken(lexer.LEFT_PAREN, "Expect ( after 'for'")

	var initializer ast.Stmt
	if p.match(lexer.SEMICOLON) {
		initializer = nil
	} else if p.match(lexer.VAR) {
		initializer = p.varDeclaration()
	} else if p.match(lexer.EQUAL) {
		initializer = p.expressionStatement()
	}

	var condition ast.Expr
	if !p.check(lexer.SEMICOLON) {
		condition = p.expression()
	}
	p.acceptToken(lexer.SEMICOLON, "Expect ; after loop condition")

	var increment ast.Expr
	if !p.check(lexer.RIGHT_PAREN) {
		increment = p.expression()
	}
	p.acceptToken(lexer.RIGHT_PAREN, "Expect ) after for clauses")

	body := p.statement()
	if increment != nil {
    body = ast.Block{Statements: []ast.Stmt{body, ast.Expression{Expression: increment}}}
	}
	if condition == nil {
    condition = ast.Literal{Value: true}
	}
  body = ast.While{Condition: condition, Body: body}

	if initializer != nil {
    body = ast.Block{Statements: []ast.Stmt{initializer, body}}
	}

	return body
}

func (p *Parser) returnStatement() ast.Stmt {
	keyword := p.previous()

	var value ast.Expr
	if !p.check(lexer.SEMICOLON) {
		value = p.expression()
	}
	p.acceptToken(lexer.SEMICOLON, "Expect ';' after return value")

  return ast.Return{Keyword: *keyword, Value: value}
}

func (p *Parser) expression() ast.Expr {
	return p.assignment()
}

func (p *Parser) assignment() ast.Expr {
	expr := p.or()

	if p.match(lexer.EQUAL) {
		equals := p.previous()
		value := p.equality()

		switch expr.(type) {
		case ast.Variable:
			name := expr.(ast.Variable).Name
      return ast.Assign{Name: name, Value: value}
		}

		p.error(equals, errors.New("Invalid assignment target"))
	}

	return expr
}

func (p *Parser) or() ast.Expr {
	expr := p.and()

	for p.match(lexer.OR) {
		operator := p.previous()
		right := p.and()
    expr = ast.Logical{Left: expr, Operator: *operator, Right: right}
	}

	return expr
}

func (p *Parser) and() ast.Expr {
	expr := p.equality()

	for p.match(lexer.AND) {
		operator := p.previous()
		right := p.equality()
    expr = ast.Logical{Left: expr, Operator: *operator, Right: right}
	}

	return expr
}

func (p *Parser) equality() ast.Expr {
	expr := p.comprasion()

	for p.match(lexer.BANG_EQUAL, lexer.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comprasion()

    expr = ast.Binary{Left: expr, Operator: *operator, Right: right}
	}

	return expr
}

func (p *Parser) comprasion() ast.Expr {
	expr := p.term()

	for p.match(lexer.GREATER, lexer.GREATER_EQUAL, lexer.LESS, lexer.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()

    expr = ast.Binary{Left: expr, Operator: *operator, Right: right}
	}

	return expr
}

func (p *Parser) term() ast.Expr {
	expr := p.factor()

	for p.match(lexer.MINUS, lexer.PLUS) {
		operator := p.previous()
		right := p.factor()

    expr = ast.Binary{Left: expr, Operator: *operator, Right: right}
	}

	return expr
}

func (p *Parser) factor() ast.Expr {
	expr := p.unary()

	for p.match(lexer.STAR, lexer.SLASH) {
		operator := p.previous()
		right := p.unary()

    expr = ast.Binary{Left: expr, Operator: *operator, Right: right}
	}

	return expr
}

func (p *Parser) unary() ast.Expr {
	if p.match(lexer.BANG, lexer.MINUS) {
		operator := p.previous()
		right := p.unary()

    return ast.Unary{Operator: *operator, Right: right}
	}

	return p.call()
}

func (p *Parser) call() ast.Expr {
	expr := p.primary()

	for {
		if p.match(lexer.LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) finishCall(callee ast.Expr) ast.Expr {
	var arguments []ast.Expr

	if !p.check(lexer.RIGHT_PAREN) {
		arguments = append(arguments, p.expression())
	}

	for p.match(lexer.COMMA) {
		arguments = append(arguments, p.expression())
		if len(arguments) > 255 {
			p.error(p.peek(), errors.New("Can't have more than 255 arguments"))
		}

		if p.check(lexer.RIGHT_PAREN) {
			break
		}

	}
	paren := p.acceptToken(lexer.RIGHT_PAREN, "Expect ')' after arguments")

  return ast.Call{Callee: callee, Paren: *paren, Arguments: arguments}
}

func (p *Parser) primary() ast.Expr {
	if p.match(lexer.TRUE) {
    return ast.Literal{Value: true}
	}
	if p.match(lexer.FALSE) {
    return ast.Literal{Value: false}
	}
	if p.match(lexer.NIL) {
    return ast.Literal{Value: nil}
	}

	if p.match(lexer.NUMBER, lexer.STRING) {
    return ast.Literal{Value: p.previous().Literal}
	}
	if p.match(lexer.IDENTIFIER) {
    return ast.Variable{Name: *p.previous()}
	}

	if p.match(lexer.LEFT_PAREN) {
		expr := p.expression()
		p.acceptToken(lexer.RIGHT_PAREN, "Expect ')' after expression")
		return ast.Grouping{Expr: expr}
	}

	p.parseError("Expect expression")
	return nil
}

func (p *Parser) errorRecovery() {
	if err := recover(); err != nil {
		p.error(p.peek(), fmt.Errorf("%v", err))
		p.synchronize()
	}
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.eof() {
		if p.previous().TokenType == lexer.SEMICOLON {
			return
		}
		switch p.peek().TokenType {
		case lexer.IF, lexer.VAR, lexer.FUN, lexer.FOR,
			lexer.WHILE, lexer.PRINT, lexer.CLASS, lexer.RETURN:
		}
		p.advance()
	}
}
