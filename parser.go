package main

import (
	"fmt"
	"strings"
)

// LatexNode is the interface for all AST nodes
type LatexNode interface {
	NodeType() string
}

// NumberNode represents a numeric literal
type NumberNode struct {
	Value float64
	Token Token
}

func (n *NumberNode) NodeType() string { return "NumberNode" }

// BinaryOpNode represents a binary operation
type BinaryOpNode struct {
	Left     LatexNode
	Operator Token
	Right    LatexNode
}

func (n *BinaryOpNode) NodeType() string { return "BinaryOpNode" }

// GroupNode represents a grouped expression (parentheses)
type GroupNode struct {
	Inner LatexNode
	Token Token
}

func (n *GroupNode) NodeType() string { return "GroupNode" }

// Parser parses tokens into a LaTeX AST
type Parser struct {
	lexer        *Lexer
	currentToken Token
	peekToken    Token
	errors       []string
}

// Precedence levels for operators
const (
	_ int = iota
	LOWEST
	SUM     // +, -
	PRODUCT // *, /
)

// precedences maps token types to their precedence
var precedences = map[TokenType]int{
	PLUS:     SUM,
	MINUS:    SUM,
	MULTIPLY: PRODUCT,
	DIVIDE:   PRODUCT,
}

// NewParser creates a new Parser instance
func NewParser(lexer *Lexer) *Parser {
	p := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

	// Read two tokens to initialize current and peek
	p.nextToken()
	p.nextToken()

	return p
}

// nextToken advances to the next token
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// Parse parses the input and returns the root AST node
func (p *Parser) Parse() (LatexNode, error) {
	node := p.parseExpression(LOWEST)

	if len(p.errors) > 0 {
		return nil, fmt.Errorf("parse errors: %s", strings.Join(p.errors, "; "))
	}

	return node, nil
}

// parseExpression implements Pratt parsing
func (p *Parser) parseExpression(precedence int) LatexNode {
	// Parse prefix expression
	var left LatexNode

	switch p.currentToken.Type {
	case NUMBER:
		left = p.parseNumber()
	case LPAREN:
		left = p.parseGroupExpression()
	default:
		p.errors = append(p.errors, fmt.Sprintf("unexpected token at position %d: %s", p.currentToken.Pos, p.currentToken.Literal))
		return nil
	}

	// Parse infix expressions with precedence climbing
	for p.peekToken.Type != EOF && precedence < p.peekPrecedence() {
		switch p.peekToken.Type {
		case PLUS, MINUS, MULTIPLY, DIVIDE:
			p.nextToken()
			left = p.parseBinaryOp(left)
		default:
			return left
		}
	}

	return left
}

// parseNumber parses a number literal
func (p *Parser) parseNumber() LatexNode {
	return &NumberNode{
		Value: p.currentToken.Value,
		Token: p.currentToken,
	}
}

// parseGroupExpression parses a parenthesized expression
func (p *Parser) parseGroupExpression() LatexNode {
	token := p.currentToken // Save the '(' token

	p.nextToken() // Move past '('

	inner := p.parseExpression(LOWEST)

	if !p.expectPeek(RPAREN) {
		p.errors = append(p.errors, fmt.Sprintf("expected ')' at position %d", p.peekToken.Pos))
		return nil
	}

	return &GroupNode{
		Inner: inner,
		Token: token,
	}
}

// parseBinaryOp parses a binary operation
func (p *Parser) parseBinaryOp(left LatexNode) LatexNode {
	node := &BinaryOpNode{
		Left:     left,
		Operator: p.currentToken,
	}

	precedence := p.currentPrecedence()
	p.nextToken()
	node.Right = p.parseExpression(precedence)

	return node
}

// expectPeek checks if the next token is of the expected type
func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	return false
}

// currentPrecedence returns the precedence of the current token
func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}

// peekPrecedence returns the precedence of the peek token
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// Errors returns the list of parse errors
func (p *Parser) Errors() []string {
	return p.errors
}

// ParseLatex parses a LaTeX string and returns an Expression tree
func ParseLatex(input string) (Expression, error) {
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	ast, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	converter := NewConverter()
	expr, err := converter.Convert(ast)
	if err != nil {
		return nil, fmt.Errorf("conversion error: %w", err)
	}

	return expr, nil
}

// ParseAndEval parses a LaTeX string and evaluates it, returning the result
func ParseAndEval(input string) (*NumberValue, error) {
	expr, err := ParseLatex(input)
	if err != nil {
		return nil, err
	}

	result, ok := expr.Eval()
	if !ok {
		return nil, fmt.Errorf("evaluation failed")
	}

	num, ok := result.(*NumberValue)
	if !ok {
		return nil, fmt.Errorf("result is not a number")
	}

	return num, nil
}
