package main

import "testing"

func TestLexer_Numbers(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"123", 123.0},
		{"3.14", 3.14},
		{"0.5", 0.5},
		{"42", 42.0},
		{"0", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok := lexer.NextToken()

			if tok.Type != NUMBER {
				t.Errorf("expected NUMBER token, got %v", tok.Type)
			}
			if tok.Value != tt.expected {
				t.Errorf("expected value %f, got %f", tt.expected, tok.Value)
			}
			if tok.Literal != tt.input {
				t.Errorf("expected literal %s, got %s", tt.input, tok.Literal)
			}
		})
	}
}

func TestLexer_Operators(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"+", PLUS},
		{"-", MINUS},
		{"*", MULTIPLY},
		{"/", DIVIDE},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok := lexer.NextToken()

			if tok.Type != tt.expected {
				t.Errorf("expected token type %v, got %v", tt.expected, tok.Type)
			}
			if tok.Literal != tt.input {
				t.Errorf("expected literal %s, got %s", tt.input, tok.Literal)
			}
		})
	}
}

func TestLexer_Parentheses(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"(", LPAREN},
		{")", RPAREN},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok := lexer.NextToken()

			if tok.Type != tt.expected {
				t.Errorf("expected token type %v, got %v", tt.expected, tok.Type)
			}
		})
	}
}

func TestLexer_Expression(t *testing.T) {
	input := "2 + 3 * 4"
	lexer := NewLexer(input)

	expectedTokens := []struct {
		tokenType TokenType
		literal   string
		value     float64
	}{
		{NUMBER, "2", 2.0},
		{PLUS, "+", 0},
		{NUMBER, "3", 3.0},
		{MULTIPLY, "*", 0},
		{NUMBER, "4", 4.0},
		{EOF, "", 0},
	}

	for i, expected := range expectedTokens {
		tok := lexer.NextToken()

		if tok.Type != expected.tokenType {
			t.Errorf("token[%d] - expected type %v, got %v", i, expected.tokenType, tok.Type)
		}
		if tok.Literal != expected.literal {
			t.Errorf("token[%d] - expected literal %s, got %s", i, expected.literal, tok.Literal)
		}
		if expected.tokenType == NUMBER && tok.Value != expected.value {
			t.Errorf("token[%d] - expected value %f, got %f", i, expected.value, tok.Value)
		}
	}
}

func TestLexer_Whitespace(t *testing.T) {
	input := "  2  +  3  "
	lexer := NewLexer(input)

	expectedTokens := []TokenType{NUMBER, PLUS, NUMBER, EOF}

	for i, expected := range expectedTokens {
		tok := lexer.NextToken()
		if tok.Type != expected {
			t.Errorf("token[%d] - expected type %v, got %v", i, expected, tok.Type)
		}
	}
}

func TestLexer_ComplexExpression(t *testing.T) {
	input := "(2 + 3) * 4 / 2"
	lexer := NewLexer(input)

	expectedTokens := []TokenType{
		LPAREN, NUMBER, PLUS, NUMBER, RPAREN,
		MULTIPLY, NUMBER, DIVIDE, NUMBER, EOF,
	}

	for i, expected := range expectedTokens {
		tok := lexer.NextToken()
		if tok.Type != expected {
			t.Errorf("token[%d] - expected type %v, got %v", i, expected, tok.Type)
		}
	}
}

func TestLexer_EmptyInput(t *testing.T) {
	lexer := NewLexer("")
	tok := lexer.NextToken()

	if tok.Type != EOF {
		t.Errorf("expected EOF for empty input, got %v", tok.Type)
	}
}

func TestLexer_IllegalCharacter(t *testing.T) {
	input := "2 @ 3"
	lexer := NewLexer(input)

	tokens := []Token{}
	for {
		tok := lexer.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == EOF {
			break
		}
	}

	if tokens[1].Type != ILLEGAL {
		t.Errorf("expected ILLEGAL token for '@', got %v", tokens[1].Type)
	}
}

func TestLexer_DecimalNumbers(t *testing.T) {
	input := "2.5 + 1.5"
	lexer := NewLexer(input)

	tok1 := lexer.NextToken()
	if tok1.Type != NUMBER || tok1.Value != 2.5 {
		t.Errorf("expected NUMBER 2.5, got %v with value %f", tok1.Type, tok1.Value)
	}

	tok2 := lexer.NextToken()
	if tok2.Type != PLUS {
		t.Errorf("expected PLUS, got %v", tok2.Type)
	}

	tok3 := lexer.NextToken()
	if tok3.Type != NUMBER || tok3.Value != 1.5 {
		t.Errorf("expected NUMBER 1.5, got %v with value %f", tok3.Type, tok3.Value)
	}
}

func TestLexer_Position(t *testing.T) {
	input := "2+3"
	lexer := NewLexer(input)

	tok1 := lexer.NextToken()
	if tok1.Pos != 0 {
		t.Errorf("expected position 0 for first token, got %d", tok1.Pos)
	}

	tok2 := lexer.NextToken()
	if tok2.Pos != 1 {
		t.Errorf("expected position 1 for second token, got %d", tok2.Pos)
	}

	tok3 := lexer.NextToken()
	if tok3.Pos != 2 {
		t.Errorf("expected position 2 for third token, got %d", tok3.Pos)
	}
}
