package main

import "testing"

func TestParser_Number(t *testing.T) {
	input := "42"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	node, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	numNode, ok := node.(*NumberNode)
	if !ok {
		t.Fatalf("expected NumberNode, got %T", node)
	}

	if numNode.Value != 42.0 {
		t.Errorf("expected value 42.0, got %f", numNode.Value)
	}
}

func TestParser_Addition(t *testing.T) {
	input := "2 + 3"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	node, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	binOp, ok := node.(*BinaryOpNode)
	if !ok {
		t.Fatalf("expected BinaryOpNode, got %T", node)
	}

	if binOp.Operator.Type != PLUS {
		t.Errorf("expected PLUS operator, got %v", binOp.Operator.Type)
	}

	left, ok := binOp.Left.(*NumberNode)
	if !ok || left.Value != 2.0 {
		t.Errorf("expected left operand 2.0")
	}

	right, ok := binOp.Right.(*NumberNode)
	if !ok || right.Value != 3.0 {
		t.Errorf("expected right operand 3.0")
	}
}

func TestParser_Multiplication(t *testing.T) {
	input := "2 * 3"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	node, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	binOp, ok := node.(*BinaryOpNode)
	if !ok {
		t.Fatalf("expected BinaryOpNode, got %T", node)
	}

	if binOp.Operator.Type != MULTIPLY {
		t.Errorf("expected MULTIPLY operator, got %v", binOp.Operator.Type)
	}
}

func TestParser_Precedence(t *testing.T) {
	input := "2 + 3 * 4"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	node, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// Should parse as: 2 + (3 * 4)
	binOp, ok := node.(*BinaryOpNode)
	if !ok {
		t.Fatalf("expected BinaryOpNode at root, got %T", node)
	}

	if binOp.Operator.Type != PLUS {
		t.Errorf("expected PLUS at root, got %v", binOp.Operator.Type)
	}

	// Left should be 2
	left, ok := binOp.Left.(*NumberNode)
	if !ok || left.Value != 2.0 {
		t.Errorf("expected left operand 2.0")
	}

	// Right should be (3 * 4)
	rightBinOp, ok := binOp.Right.(*BinaryOpNode)
	if !ok {
		t.Fatalf("expected BinaryOpNode on right, got %T", binOp.Right)
	}

	if rightBinOp.Operator.Type != MULTIPLY {
		t.Errorf("expected MULTIPLY on right, got %v", rightBinOp.Operator.Type)
	}

	rightLeft, ok := rightBinOp.Left.(*NumberNode)
	if !ok || rightLeft.Value != 3.0 {
		t.Errorf("expected right-left operand 3.0")
	}

	rightRight, ok := rightBinOp.Right.(*NumberNode)
	if !ok || rightRight.Value != 4.0 {
		t.Errorf("expected right-right operand 4.0")
	}
}

func TestParser_LeftAssociativity(t *testing.T) {
	input := "10 - 3 - 2"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	node, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// Should parse as: (10 - 3) - 2
	binOp, ok := node.(*BinaryOpNode)
	if !ok {
		t.Fatalf("expected BinaryOpNode at root, got %T", node)
	}

	if binOp.Operator.Type != MINUS {
		t.Errorf("expected MINUS at root, got %v", binOp.Operator.Type)
	}

	// Left should be (10 - 3)
	leftBinOp, ok := binOp.Left.(*BinaryOpNode)
	if !ok {
		t.Fatalf("expected BinaryOpNode on left, got %T", binOp.Left)
	}

	if leftBinOp.Operator.Type != MINUS {
		t.Errorf("expected MINUS on left, got %v", leftBinOp.Operator.Type)
	}

	// Right should be 2
	right, ok := binOp.Right.(*NumberNode)
	if !ok || right.Value != 2.0 {
		t.Errorf("expected right operand 2.0")
	}
}

func TestParser_Grouping(t *testing.T) {
	input := "(2 + 3) * 4"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	node, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// Should parse as: (2 + 3) * 4
	binOp, ok := node.(*BinaryOpNode)
	if !ok {
		t.Fatalf("expected BinaryOpNode at root, got %T", node)
	}

	if binOp.Operator.Type != MULTIPLY {
		t.Errorf("expected MULTIPLY at root, got %v", binOp.Operator.Type)
	}

	// Left should be GroupNode containing (2 + 3)
	groupNode, ok := binOp.Left.(*GroupNode)
	if !ok {
		t.Fatalf("expected GroupNode on left, got %T", binOp.Left)
	}

	innerBinOp, ok := groupNode.Inner.(*BinaryOpNode)
	if !ok {
		t.Fatalf("expected BinaryOpNode inside group, got %T", groupNode.Inner)
	}

	if innerBinOp.Operator.Type != PLUS {
		t.Errorf("expected PLUS inside group, got %v", innerBinOp.Operator.Type)
	}

	// Right should be 4
	right, ok := binOp.Right.(*NumberNode)
	if !ok || right.Value != 4.0 {
		t.Errorf("expected right operand 4.0")
	}
}

func TestParser_NestedGroups(t *testing.T) {
	input := "((2 + 3))"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	node, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// Outer group
	outerGroup, ok := node.(*GroupNode)
	if !ok {
		t.Fatalf("expected GroupNode at root, got %T", node)
	}

	// Inner group
	innerGroup, ok := outerGroup.Inner.(*GroupNode)
	if !ok {
		t.Fatalf("expected GroupNode inside outer group, got %T", outerGroup.Inner)
	}

	// 2 + 3
	binOp, ok := innerGroup.Inner.(*BinaryOpNode)
	if !ok {
		t.Fatalf("expected BinaryOpNode inside inner group, got %T", innerGroup.Inner)
	}

	if binOp.Operator.Type != PLUS {
		t.Errorf("expected PLUS, got %v", binOp.Operator.Type)
	}
}

func TestParser_ErrorUnexpectedEOF(t *testing.T) {
	input := "2 +"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	_, err := parser.Parse()
	if err == nil {
		t.Errorf("expected error for incomplete expression")
	}
}

func TestParser_ErrorUnmatchedParen(t *testing.T) {
	input := "(2 + 3"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	_, err := parser.Parse()
	if err == nil {
		t.Errorf("expected error for unmatched parenthesis")
	}
}

func TestParser_ErrorInvalidSyntax(t *testing.T) {
	input := "2 + + 3"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	_, err := parser.Parse()
	if err == nil {
		t.Errorf("expected error for invalid syntax")
	}
}

func TestParser_ComplexExpression(t *testing.T) {
	input := "(1 + 2) * (3 + 4)"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	node, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	binOp, ok := node.(*BinaryOpNode)
	if !ok {
		t.Fatalf("expected BinaryOpNode at root, got %T", node)
	}

	if binOp.Operator.Type != MULTIPLY {
		t.Errorf("expected MULTIPLY at root, got %v", binOp.Operator.Type)
	}

	_, ok = binOp.Left.(*GroupNode)
	if !ok {
		t.Errorf("expected GroupNode on left")
	}

	_, ok = binOp.Right.(*GroupNode)
	if !ok {
		t.Errorf("expected GroupNode on right")
	}
}
