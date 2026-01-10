package latex

import (
	"exprtree/expr"
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

// VariableNode represents a variable
type VariableNode struct {
	Name  string
	Token Token
}

func (n *VariableNode) NodeType() string { return "VariableNode" }

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

// CommandNode represents a LaTeX command like \sqrt
type CommandNode struct {
	Name     string
	Argument LatexNode
	Optional LatexNode // nil if not present
	Token    Token
}

func (n *CommandNode) NodeType() string { return "CommandNode" }

// EqualNode represents an equality expression
type EqualNode struct {
	Left     LatexNode
	Operator Token
	Right    LatexNode
}

func (n *EqualNode) NodeType() string { return "EqualNode" }

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
	EQUALITY // =
	SUM      // +, -
	PRODUCT  // *, /
	POWER    // ^
)

// precedences maps token types to their precedence
var precedences = map[TokenType]int{
	EQUAL:    EQUALITY,
	PLUS:     SUM,
	MINUS:    SUM,
	MULTIPLY: PRODUCT,
	DIVIDE:   PRODUCT,
	CARET:    POWER,
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
	case VARIABLE:
		left = p.parseVariable()
	case LPAREN:
		left = p.parseGroupExpression()
	case COMMAND:
		left = p.parseCommand()
	default:
		p.errors = append(p.errors, fmt.Sprintf("unexpected token at position %d: %s", p.currentToken.Pos, p.currentToken.Literal))
		return nil
	}

	// Parse infix expressions with precedence climbing
	for p.peekToken.Type != EOF && precedence < p.peekPrecedence() {
		switch p.peekToken.Type {
		case PLUS, MINUS, MULTIPLY, DIVIDE, CARET, EQUAL:
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

// parseVariable parses a variable
func (p *Parser) parseVariable() LatexNode {
	return &VariableNode{
		Name:  p.currentToken.Literal,
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

// parseCommand parses \sqrt{...} and \sqrt[n]{...}
func (p *Parser) parseCommand() LatexNode {
	token := p.currentToken
	commandName := token.Literal

	var optional LatexNode

	// Check for optional argument [n]
	if p.peekToken.Type == LBRACKET {
		p.nextToken() // consume LBRACKET
		p.nextToken() // move to content
		optional = p.parseExpression(LOWEST)

		if !p.expectPeek(RBRACKET) {
			p.errors = append(p.errors, fmt.Sprintf("expected ']' at position %d", p.peekToken.Pos))
			return nil
		}
	}

	// Parse required argument {expr}
	if !p.expectPeek(LBRACE) {
		p.errors = append(p.errors, fmt.Sprintf("expected '{' after \\%s at position %d", commandName, p.peekToken.Pos))
		return nil
	}

	p.nextToken() // move past LBRACE
	argument := p.parseExpression(LOWEST)

	if !p.expectPeek(RBRACE) {
		p.errors = append(p.errors, fmt.Sprintf("expected '}' at position %d", p.peekToken.Pos))
		return nil
	}

	return &CommandNode{
		Name:     commandName,
		Argument: argument,
		Optional: optional,
		Token:    token,
	}
}

// parseBinaryOp parses a binary operation
func (p *Parser) parseBinaryOp(left LatexNode) LatexNode {
	operator := p.currentToken
	precedence := p.currentPrecedence()
	p.nextToken()

	var right LatexNode
	// Right-associative for power (^), left-associative for others
	if operator.Type == CARET {
		right = p.parseExpression(precedence - 1) // Right-associative: use precedence - 1
	} else {
		right = p.parseExpression(precedence) // Left-associative: use same precedence
	}

	// Create EqualNode for equality, BinaryOpNode for other operators
	if operator.Type == EQUAL {
		return &EqualNode{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return &BinaryOpNode{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
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
func ParseLatex(input string) (expr.Expression, error) {
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	ast, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	converter := NewConverter()
	expression, err := converter.Convert(ast)
	if err != nil {
		return nil, fmt.Errorf("conversion error: %w", err)
	}

	return expression, nil
}

// ParseAndEval parses a LaTeX string and evaluates it, returning the result
func ParseAndEval(input string) (*expr.NumberValue, error) {
	expression, err := ParseLatex(input)
	if err != nil {
		return nil, err
	}

	result, ok := expression.Eval()
	if !ok {
		return nil, fmt.Errorf("evaluation failed")
	}

	num, ok := result.(*expr.NumberValue)
	if !ok {
		return nil, fmt.Errorf("result is not a number")
	}

	return num, nil
}
